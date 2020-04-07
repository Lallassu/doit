package main

import (
	"crypto/rand"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// DB holds the database information
type DB struct {
	db *gorm.DB
}

func (d *DB) HashPass(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return string(hash)
}

func (d *DB) VerifyPass(user, pass string) bool {
	// Get user
	acc := Account{}
	d.db.Unscoped().Where("user = ?", user).First(&acc)
	if acc.ID == 0 {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(acc.Pass), []byte(pass))
	if err != nil {
		return false
	}

	return true
}

func (d *DB) CreateUser(user, pass, email string, admin bool) bool {

	existing := Account{}
	d.db.Where("user = ?", user).First(&existing)
	if existing.ID != 0 {
		return false
	}

	hpass := d.HashPass(pass)
	acc := &Account{Email: email, User: user, Pass: hpass, Admin: admin}
	d.db.Create(acc)

	return true
}

func (d *DB) HasAdminAccount() bool {
	acc := Account{}
	d.db.Where("admin = ?", true).First(&acc)
	if acc.ID != 0 {
		return true
	}
	return false
}

// Init setups the database and creates tables if needed.
func (d *DB) Init(dbname string) {
	db, err := gorm.Open("sqlite3", dbname)
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&List{}, &Item{}, &Account{}, &Token{})
	d.db = db
}

func (d *DB) GetAccounts() map[string]string {
	accounts := make(map[string]string, 0)
	acc := []Account{}
	d.db.Unscoped().Find(&acc)

	for _, a := range acc {
		accounts[a.User] = a.Pass
	}
	return accounts
}

func (d *DB) ValidateToken(token string) (*Account, bool) {
	tok := Token{}
	d.db.Preload("Account").Where("token = ?", token).First(&tok)
	if tok.ID != 0 {
		return &tok.Account, true
	}
	return nil, false
}

func (d *DB) ValidateLogin(user, pass string) (*Token, bool) {
	a := Account{}
	d.db.Where("user = ?", user).First(&a)

	// Not validated yet -> fail
	if !a.Admin && !a.Validated {
		return nil, false
	}

	if !d.VerifyPass(user, pass) {
		return nil, false
	}

	b := make([]byte, 4)
	rand.Read(b)
	str := fmt.Sprintf("%x%x", b, a.User)

	tok := &Token{Account: a, Token: str}
	d.db.Save(tok)

	return tok, true
}

func (d *DB) HaveAccess(accId, listId uint) bool {
	list := &List{}
	d.db.Preload("Account").Preload("Share").First(list, listId)

	haveAccess := false
	if list.Account.ID == accId {
		haveAccess = true
	}
	for _, a := range list.Share {
		if a.ID == accId {
			haveAccess = true
			break
		}
	}
	return haveAccess
}

func (d *DB) AdminValidateAccount(id int, acc *Account) {
	if acc.Admin {
		d.db.Model(&Account{}).Where("id = ?", id).Update("validated", true)
	}
}

func (d *DB) AdminRemoveAccount(id int, acc *Account) {
	if acc.Admin && acc.ID != uint(id) {
		d.db.Unscoped().Delete(&Account{}, id)
	}
}

func (d *DB) AdminRemoveToken(id int, acc *Account) {
	if acc.Admin {
		d.db.Unscoped().Delete(&Token{}, id)
	}
}

func (d *DB) AdminAllAccounts(acc *Account) []Account {
	accounts := []Account{}
	if acc.Admin {
		d.db.Find(&accounts)
	}
	return accounts
}

func (d *DB) AdminAllTokens(acc *Account) []Token {
	tokens := []Token{}
	if acc.Admin {
		d.db.Preload("Account").Find(&tokens)
	}
	return tokens
}

func (d *DB) ItemsForList(id int, acc *Account) []Item {
	items := []Item{}
	if d.HaveAccess(acc.ID, uint(id)) {
		d.db.Unscoped().Preload("Account").Where("list_id = ?", id).Find(&items)
	}
	return items
}

func (d *DB) AllLists(acc *Account) []List {
	allLists := []List{}
	lists := []List{}

	d.db.Preload("Account").Preload("Share").Find(&allLists)
	for _, l := range allLists {
		if d.HaveAccess(acc.ID, l.ID) {
			lists = append(lists, l)
			continue
		}
	}

	return lists
}

func (d *DB) UpdateSortOrder(l []Item) {
	for _, item := range l {
		d.db.Unscoped().Model(&item).Update("order", item.Order)
	}
}

func (d *DB) LastUpdate(id int) time.Time {
	item := &Item{}
	d.db.Unscoped().Model(&List{}).Order("updated_at desc").First(item)
	return item.UpdatedAt
}

func (d *DB) RenameList(l List, acc *Account) {
	if d.HaveAccess(acc.ID, uint(l.ID)) {
		d.db.Unscoped().Model(&l).Where("account_id = ?", acc.ID).Update("Name", l.Name)
	}
}

func (d *DB) SaveList(l *List, acc *Account) {
	l.DeletedAt = nil
	l.Account = *acc
	d.db.Unscoped().Create(l)
}

func (d *DB) SaveItem(l *Item, acc *Account) {
	if d.HaveAccess(acc.ID, uint(l.ListId)) {
		l.DeletedAt = nil
		l.ReminderSent = false
		l.Account = *acc
		d.db.Unscoped().Create(l)
	}
}

func (d *DB) Favorite(fav bool, listId int, acc *Account) {
	if fav {
		acc.Favorite = listId
	} else {
		acc.Favorite = 0
	}
	d.db.Unscoped().Save(acc)
}

func (d *DB) UpdateItemSent(l *Item, acc *Account) {
	// TBD: Should not be required...?
	l.AccountID = acc.ID
	l.Account = *acc
	d.db.Unscoped().Save(l)
}

func (d *DB) UpdateItem(l *Item, acc *Account) {
	if d.HaveAccess(acc.ID, uint(l.ListId)) {
		l.DeletedAt = nil
		item := &Item{}
		d.db.Preload("Account").First(&item, l.ID)

		// Keep creator account, this is a workaround for now.
		l.Account = item.Account

		// Only set remindersent to false if time has changed.
		if l.Time != item.Time {
			l.ReminderSent = false
		}
		d.db.Unscoped().Save(l)
	}
}

func (d *DB) Logout(token string, acc *Account) {
	tok := &Token{}
	d.db.Preload("Account").Where("token = ?", token).First(&tok)

	if tok.Account.ID == acc.ID {
		d.db.Unscoped().Where("token = ?", token).Delete(&tok)
	}
}

func (d *DB) ShareList(shareWith string, listId int, acc *Account) *Account {
	// Lookup the list
	list := &List{}
	d.db.Preload("Account").Preload("Share").First(&list, listId)
	if list.ID == 0 {
		return nil
	}
	// Look up co-owner
	co := Account{}
	d.db.Where("user = ? OR email = ?", shareWith, shareWith).First(&co)
	if co.ID == 0 {
		return nil
	}

	// Check if user is already in shared list.
	for _, a := range list.Share {
		if a.ID == co.ID {
			return nil
		}
	}

	// Verify owner of list.
	if acc.ID == list.Account.ID {
		// add user to co-owner
		if co.ID != list.Account.ID {
			d.db.Model(list).Association("Share").Append(co)
			return &co
		}
	}
	return nil
}

func (d *DB) RemoveShareList(shareWith string, listId int, acc *Account) *Account {
	// Lookup the list
	list := &List{}
	d.db.Preload("Share").Preload("Account").First(&list, listId)
	if list.ID == 0 {
		return nil
	}

	// Look up co-owner
	co := Account{}
	d.db.Where("user = ? OR email = ?", shareWith, shareWith).First(&co)

	// Verify owner of list.
	if acc.ID == list.Account.ID || acc.ID == co.ID {
		d.db.Model(list).Association("Share").Delete(&co)
		return &co
	}

	return nil
}

func (d *DB) DeleteList(l *List, acc *Account) {
	list := &List{}
	d.db.Preload("Account").First(&list, l.ID)

	if list.Account.ID == acc.ID {
		d.db.Unscoped().Where("list_id = ?", list.ID).Delete(&Item{})
		d.db.Unscoped().Where("id = ?", list.ID).Delete(list)
		d.db.Model(&Account{}).Where("favorite = ?", list.ID).UpdateColumn("favorite", 0)
	}
}

func (d *DB) AllListsWithItems() []List {
	lists := []List{}

	d.db.Preload("Account").Preload("Share").Preload("Items").Find(&lists)

	return lists
}

func (d *DB) GetAccount(id uint) *Account {
	account := &Account{}
	d.db.Find(account, id)
	return account
}
