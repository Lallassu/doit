<template>
  <div class="container">
    <div v-if="errorText">
      <span class="error">Something went wrong: {{ errorText }}<br>Please reload page.</span>
    </div>

    <div v-else>
      <!-- Admin View -->
      <div v-if="account.Admin">
        <Button label="Logout" severity="danger" size="small" @click="logout()" />
        <span class="admin-text">You are admin.</span>

        <DataTable :value="accounts" stripedRows class="mt-3">
          <template #header>
            <h3>Accounts</h3>
          </template>
          <Column field="ID" header="ID" />
          <Column field="User" header="User" />
          <Column field="Email" header="Email" />
          <Column field="Admin" header="Admin" />
          <Column field="Validated" header="Validated" />
          <Column header="Actions">
            <template #body="{ data }">
              <Button v-if="!data.Validated" label="Validate" severity="success" size="small" @click="validateAccount(data.ID)" class="mr-1" />
              <Button v-if="!data.Admin" label="Remove" severity="danger" size="small" @click="removeAccount(data.ID)" />
            </template>
          </Column>
        </DataTable>

        <DataTable :value="tokens" stripedRows class="mt-3">
          <template #header>
            <h3>Tokens</h3>
          </template>
          <Column field="ID" header="ID" />
          <Column field="Account.User" header="Account" />
          <Column field="Token" header="Token" />
          <Column header="Actions">
            <template #body="{ data }">
              <Button label="Revoke" severity="danger" size="small" @click="revokeToken(data.ID)" />
            </template>
          </Column>
        </DataTable>
      </div>

      <!-- Signup Form -->
      <div v-if="createadm || createNewAcc" class="auth-form">
        <span v-if="createadm" class="admin-text">This will create the admin account!</span>
        <div class="form-field">
          <i class="fas fa-user"></i>
          <InputText v-model="username" placeholder="Username" />
        </div>
        <div class="form-field">
          <i class="fas fa-envelope"></i>
          <InputText v-model="email" placeholder="email@example.com" />
        </div>
        <div class="form-field">
          <i class="fas fa-key"></i>
          <InputText v-model="password" type="password" placeholder="Password" />
        </div>
        <div class="form-field">
          <i class="fas fa-key"></i>
          <InputText v-model="password2" type="password" placeholder="Repeat password..." @keyup.enter="signup()" />
        </div>
        <Button label="Create Account" @click="signup()" class="w-full" />
      </div>

      <!-- Login Form -->
      <div v-if="!loggedIn && !createadm && !createNewAcc" class="auth-form">
        <div class="form-field">
          <i class="fas fa-user"></i>
          <InputText v-model="username" placeholder="User" @keyup.enter="login()" />
        </div>
        <div class="form-field">
          <i class="fas fa-key"></i>
          <InputText v-model="password" type="password" placeholder="Password" @keyup.enter="login()" />
        </div>
        <Button label="Login" @click="login()" class="w-full" />
        <Button label="Sign Up" link @click="createNewAcc = true" class="w-full mt-2" />
      </div>

      <!-- Main App View -->
      <div v-if="loggedIn && !account.Admin">
        <!-- Sidebar Menu -->
        <Button icon="pi pi-bars" @click="sidebarVisible = true" class="menu-btn" />
        <Drawer v-model:visible="sidebarVisible" header="Menu" class="sidebar-menu">
          <Button label="New List" @click="newListDialogVisible = true" class="w-full mb-3" />
          <div v-for="item in lists" :key="item.ID" class="list-menu-item" @click="selectList(item); sidebarVisible = false">
            {{ item.Name }}
            <i v-if="item.ID === account.Favorite" class="fas fa-star fav-icon"></i>
          </div>
          <div class="sidebar-footer">
            <span class="logged-in-text">Logged in as {{ account.User }}</span>
            <Button label="Logout" severity="secondary" size="small" @click="logout()" />
          </div>
        </Drawer>

        <!-- New List Dialog -->
        <Dialog v-model:visible="newListDialogVisible" header="New List Name" modal>
          <InputText v-model="newListValue" placeholder="List name" class="w-full" @keyup.enter="newList()" />
          <template #footer>
            <Button label="Cancel" severity="secondary" @click="newListDialogVisible = false" />
            <Button label="Create" @click="newList()" />
          </template>
        </Dialog>

        <!-- No Lists Message -->
        <div v-if="nolists" class="text-center mt-4">
          <Button label="Create your first list" @click="newListDialogVisible = true" />
        </div>

        <!-- Active List Header -->
        <div v-if="activelist" class="list-header">
          <span class="list-name" @click="renameListDialogVisible = true">{{ activelist.Name }}</span>
          <i class="fas fa-trash action-btn delete-btn" @click="showDeleteList = true"></i>
          <i class="fas fa-user-plus action-btn share-btn" @click="showShareList = true"></i>
          <i class="fas fa-star action-btn" :class="{ 'fav-active': activelist.Favorite }" @click="setFavorite()"></i>
        </div>

        <!-- Rename List Dialog -->
        <Dialog v-model:visible="renameListDialogVisible" header="Change List Name" modal>
          <InputText v-model="newListName" placeholder="New name" class="w-full" @keyup.enter="changeListName()" />
          <template #footer>
            <Button label="Cancel" severity="secondary" @click="renameListDialogVisible = false" />
            <Button label="Update" @click="changeListName()" />
          </template>
        </Dialog>

        <!-- Delete List Dialog -->
        <Dialog v-model:visible="showDeleteList" header="Delete List?" modal>
          <p>This will remove the list AND all items in the list, are you really sure?</p>
          <template #footer>
            <Button label="Cancel" severity="secondary" @click="showDeleteList = false" />
            <Button label="Delete" severity="danger" @click="deleteList()" />
          </template>
        </Dialog>

        <!-- Share List Dialog -->
        <Dialog v-model:visible="showShareList" header="Share List With User" modal>
          <InputText v-model="shareWithUser" placeholder="Email or username" class="w-full" @keyup.enter="shareList()" />
          <template #footer>
            <Button label="Cancel" severity="secondary" @click="showShareList = false" />
            <Button label="Share" @click="shareList()" />
          </template>
        </Dialog>

        <!-- Add Item Input -->
        <div v-if="activelist" class="add-item">
          <div class="add-icon" @click="addItem()">
            <i class="fas fa-plus"></i>
          </div>
          <InputText v-model="entry" placeholder="Add..." @keyup.enter="addItem()" class="add-input" />
        </div>

        <!-- Todo List -->
        <div class="todo-list">
          <draggable
            v-model="todoList"
            item-key="ID"
            @end="onSortEnd"
            :animation="200"
            ghost-class="drag-ghost"
            drag-class="drag-active"
          >
            <template #item="{ element: item }">
              <div class="item" :class="{ expanded: item.showNote || item.showTime }" :style="{ minHeight: item.itemSize + 'px' }">
                <div class="item-content">
                  <div class="item-checkbox" @click="completeItem(item)">
                    <i :class="item.Complete ? 'fas fa-check-square' : 'fas fa-square'"></i>
                  </div>
                  <div class="item-title" @dblclick="changeTitle(item)">
                    <span v-if="!item.editing">{{ item.Title }}</span>
                    <InputText v-else v-model="item.Title" @keyup.enter="saveTitle(item)" @keyup.esc="cancelEditTitle(item)" @blur="saveTitle(item)" class="title-input" />
                  </div>
                  <div v-if="activelist.Share && activelist.Share.length > 0" class="item-owner">
                    [{{ item.Account.User }}]
                  </div>
                  <div class="item-actions">
                    <i class="fas fa-copy action-icon" @click="copyItem(item)"></i>
                    <i class="fas fa-pen action-icon" :class="{ active: item.Note.length > 0 }" @click="showNote(item)"></i>
                    <i class="fas fa-bell action-icon" :class="{ active: item.Time.length > 0 }" @click="showTime(item)"></i>
                  </div>
                </div>

                <!-- Note Panel -->
                <div v-if="item.showNote" class="item-panel">
                  <Textarea v-model="item.Note" rows="5" class="w-full" @input="item.changed = true" />
                  <div class="panel-actions">
                    <Button label="Save" size="small" @click="saveNote(item)" />
                    <Button label="Cancel" severity="secondary" size="small" @click="showNote(item)" />
                    <span class="autosave-indicator" :class="{ visible: item.autosaved }">Auto saved.</span>
                  </div>
                </div>

                <!-- Alarm Panel -->
                <div v-if="item.showTime" class="item-panel alarm-panel">
                  <VueDatePicker v-model="item.dateTime" :enable-time-picker="true" :format="'yyyy-MM-dd HH:mm'" />
                  <div class="alarm-options">
                    <div class="alarm-option">
                      <span class="alarm-label">Notify me</span>
                      <InputNumber v-model="item.PreAlarmMinutes" :min="0" :showButtons="false" class="alarm-input" />
                      <span class="alarm-unit">minutes before</span>
                    </div>
                    <div class="alarm-option">
                      <span class="alarm-label">Repeat every</span>
                      <InputNumber v-model="item.RecurDays" :min="0" :showButtons="false" class="alarm-input" />
                      <span class="alarm-unit">days</span>
                    </div>
                    <div v-if="telegramChats.length > 0" class="alarm-option">
                      <span class="alarm-label">Telegram</span>
                      <Select v-model="item.TelegramChat" :options="telegramChats" placeholder="None" showClear class="telegram-select" />
                    </div>
                  </div>
                  <div class="panel-actions">
                    <Button label="Set Alarm" size="small" @click="saveTime(item)" />
                    <Button label="Clear" severity="secondary" size="small" @click="clearTime(item)" />
                  </div>
                </div>
              </div>
            </template>
          </draggable>
        </div>

        <!-- Completed Toggle -->
        <div v-if="completedToDoList.length > 0" class="show-completed-wrapper">
          <div class="show-completed-btn" @click="showCompletedList = !showCompletedList">
            <i :class="showCompletedList ? 'fas fa-chevron-up' : 'fas fa-chevron-down'"></i>
            {{ showCompletedList ? 'Hide' : 'Show' }} completed ({{ completedToDoList.length }})
          </div>
        </div>

        <!-- Completed List -->
        <div v-if="showCompletedList" class="todo-list completed-list">
          <div v-for="item in completedToDoList" :key="item.ID" class="item" :class="{ expanded: item.showNote || item.showTime }" :style="{ minHeight: item.itemSize + 'px' }">
            <div class="item-content">
              <div class="item-checkbox" @click="uncompleteItem(item)">
                <i class="fas fa-check-square"></i>
              </div>
              <div class="item-title">{{ item.Title }}</div>
              <div v-if="activelist.Share && activelist.Share.length > 0" class="item-owner">
                [{{ item.Account.User }}]
              </div>
              <div class="item-actions">
                <i class="fas fa-copy action-icon" @click="copyItem(item)"></i>
                <i class="fas fa-pen action-icon" :class="{ active: item.Note.length > 0 }" @click="showNote(item)"></i>
                <i class="fas fa-bell action-icon" :class="{ active: item.Time.length > 0 }" @click="showTime(item)"></i>
              </div>
            </div>

            <div v-if="item.showNote" class="item-panel">
              <Textarea v-model="item.Note" rows="5" class="w-full" />
              <div class="panel-actions">
                <Button label="Save" size="small" @click="saveNote(item)" />
                <Button label="Cancel" severity="secondary" size="small" @click="showNote(item)" />
              </div>
            </div>

            <div v-if="item.showTime" class="item-panel alarm-panel">
              <VueDatePicker v-model="item.dateTime" :enable-time-picker="true" :format="'yyyy-MM-dd HH:mm'" />
              <div class="alarm-options">
                <div class="alarm-option">
                  <span class="alarm-label">Notify me</span>
                  <InputNumber v-model="item.PreAlarmMinutes" :min="0" :showButtons="false" class="alarm-input" />
                  <span class="alarm-unit">minutes before</span>
                </div>
                <div class="alarm-option">
                  <span class="alarm-label">Repeat every</span>
                  <InputNumber v-model="item.RecurDays" :min="0" :showButtons="false" class="alarm-input" />
                  <span class="alarm-unit">days</span>
                </div>
                <div v-if="telegramChats.length > 0" class="alarm-option">
                  <span class="alarm-label">Telegram</span>
                  <Select v-model="item.TelegramChat" :options="telegramChats" placeholder="None" showClear class="telegram-select" />
                </div>
              </div>
              <div class="panel-actions">
                <Button label="Set Alarm" size="small" @click="saveTime(item)" />
                <Button label="Clear" severity="secondary" size="small" @click="clearTime(item)" />
              </div>
            </div>
          </div>
          <div class="clear-completed-wrapper">
            <Button label="Clear Completed" severity="danger" size="small" @click="clearCompleted()" />
          </div>
        </div>

        <!-- Shared Users -->
        <div v-if="activelist && activelist.Share && activelist.Share.length > 0" class="shared-users">
          <Button :label="activelist.Account.User" icon="pi pi-user" size="small" severity="secondary" />
          <Button v-for="user in activelist.Share" :key="user.ID" :label="user.User" icon="pi pi-users" size="small" severity="success" @click="confirmRemoveShare(user)" />
        </div>

        <!-- Remove Share Dialog -->
        <Dialog v-model:visible="showRemoveShareList" header="Remove Share?" modal>
          <p>This will remove the user from shared list, are you really sure?</p>
          <template #footer>
            <Button label="Cancel" severity="secondary" @click="showRemoveShareList = false" />
            <Button label="Remove" severity="danger" @click="removeUserShare()" />
          </template>
        </Dialog>
      </div>
    </div>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
import { VueDatePicker } from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Drawer from 'primevue/drawer'
import Select from 'primevue/select'

const timeFormat = "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'"

export default {
  name: 'ToDoList',
  components: {
    draggable,
    VueDatePicker,
    Button,
    InputText,
    InputNumber,
    Textarea,
    Dialog,
    DataTable,
    Column,
    Drawer,
    Select
  },
  data() {
    return {
      account: {},
      tokens: [],
      accounts: [],
      nolists: false,
      createNewAcc: false,
      email: '',
      createadm: false,
      username: '',
      password: '',
      password2: '',
      removeUser: null,
      shareWithUser: '',
      showRemoveShareList: false,
      showShareList: false,
      showDeleteList: false,
      newListDialogVisible: false,
      newListValue: '',
      loggedIn: false,
      lastUpdate: 0,
      newListName: '',
      renameListDialogVisible: false,
      activelist: null,
      itemSize: 40,
      errorText: '',
      entry: '',
      lists: [],
      todoList: [],
      completedToDoList: [],
      showCompletedList: false,
      sidebarVisible: false,
      telegramChats: []
    }
  },
  created() {
    this.fetchTelegramChats()
    if (localStorage.token) {
      this.loggedIn = true
      this.account = JSON.parse(localStorage.getItem('account') || '{}')
      this.selectList(JSON.parse(localStorage.getItem('activelist') || 'null'))

      if (this.account.Admin) {
        this.fetchAdminItems()
      } else {
        this.fetchLists()
        setInterval(() => {
          if (this.loggedIn) this.periodicalUpdate()
        }, 5000)
        setInterval(() => {
          if (this.loggedIn) this.saveActiveNotes()
        }, 5000)
      }
    } else {
      this.checkAdmin()
    }
  },
  methods: {
    async apiRequest(method, url, data = null) {
      const options = {
        method,
        headers: {
          'Content-Type': 'application/json',
          'token': localStorage.token || ''
        }
      }
      if (data) {
        options.body = JSON.stringify(data)
      }
      const response = await fetch(url, options)
      if (response.status === 401) {
        this.loggedIn = false
        localStorage.clear()
        window.location.reload()
        throw new Error('Unauthorized')
      }
      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`)
      }
      const text = await response.text()
      return text ? JSON.parse(text) : null
    },

    async checkAdmin() {
      try {
        const hasAdmin = await fetch('/hasadm').then(r => r.json())
        this.createadm = !hasAdmin
      } catch (err) {
        console.error(err)
      }
    },

    async fetchTelegramChats() {
      try {
        const chats = await fetch('/telegram-chats').then(r => r.json())
        this.telegramChats = chats || []
      } catch (err) {
        this.telegramChats = []
      }
    },

    async revokeToken(id) {
      await this.apiRequest('PUT', `/api/removetoken/${id}`)
      this.fetchAdminItems()
    },

    async validateAccount(id) {
      await this.apiRequest('PUT', `/api/validate/${id}`)
      this.fetchAdminItems()
    },

    async removeAccount(id) {
      await this.apiRequest('PUT', `/api/removeaccount/${id}`)
      this.fetchAdminItems()
    },

    async fetchAdminItems() {
      try {
        this.accounts = await this.apiRequest('GET', '/api/allaccounts')
        this.tokens = await this.apiRequest('GET', '/api/alltokens')
      } catch (err) {
        console.error(err)
      }
    },

    async signup() {
      if (!this.email || !this.password || !this.password2 || !this.username || this.password !== this.password2) {
        return
      }
      try {
        await fetch('/signup', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            user: this.username,
            pass: this.password,
            email: this.email,
            adm: this.createadm
          })
        })
        window.location.reload()
      } catch (err) {
        console.error(err)
      }
    },

    async setFavorite() {
      this.activelist.Favorite = !this.activelist.Favorite
      this.account.Favorite = this.activelist.Favorite ? this.activelist.ID : 0
      localStorage.setItem('account', JSON.stringify(this.account))

      await this.apiRequest('POST', '/api/favorite', {
        list: this.activelist.ID,
        favorite: this.activelist.Favorite
      })
      this.fetchLists()
    },

    confirmRemoveShare(user) {
      this.removeUser = user
      this.showRemoveShareList = true
    },

    async removeUserShare() {
      if (!this.removeUser) return
      await this.apiRequest('POST', '/api/removesharelist', {
        user: this.removeUser.User,
        list: this.activelist.ID
      })
      const index = this.activelist.Share.findIndex(u => u.ID === this.removeUser.ID)
      if (index > -1) this.activelist.Share.splice(index, 1)
      this.removeUser = null
      this.showRemoveShareList = false
      this.fetchLists()
    },

    async shareList() {
      const result = await this.apiRequest('POST', '/api/sharelist', {
        user: this.shareWithUser,
        list: this.activelist.ID
      })
      if (result) {
        if (!this.activelist.Share) this.activelist.Share = []
        this.activelist.Share.push(result)
      }
      this.shareWithUser = ''
      this.showShareList = false
      this.fetchLists()
    },

    async clearCompleted() {
      await this.apiRequest('POST', '/api/deletecompleted', this.activelist)
      this.selectList(this.activelist)
      this.showCompletedList = false
    },

    async deleteList() {
      if (this.account.Favorite === this.activelist.ID) {
        this.account.Favorite = 0
      }
      await this.apiRequest('POST', '/api/deletelist', this.activelist)
      this.todoList = []
      this.completedToDoList = []
      this.activelist = null
      this.showDeleteList = false
      this.fetchLists()
    },

    logout() {
      this.apiRequest('POST', '/api/logout').catch(() => {})
      this.loggedIn = false
      this.account = {}
      localStorage.clear()
    },

    async login() {
      if (!this.password || !this.username) return
      try {
        const response = await fetch('/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ user: this.username, pass: this.password })
        })
        if (!response.ok) {
          this.password = ''
          return
        }
        const data = await response.json()
        this.loggedIn = true
        localStorage.token = data.Token
        localStorage.setItem('account', JSON.stringify(data.Account))
        this.account = data.Account
        this.password = ''
        if (this.account.Admin) {
          this.fetchAdminItems()
        } else {
          this.fetchLists()
        }
      } catch (err) {
        console.error(err)
      }
    },

    async newList() {
      if (!this.newListValue) return
      const result = await this.apiRequest('POST', '/api/list', { Name: this.newListValue })
      this.activelist = result
      this.newListValue = ''
      this.newListDialogVisible = false
      this.fetchLists()
    },

    saveActiveNotes() {
      for (const item of this.todoList) {
        if (item.showNote && item.changed) {
          this.saveItem(item, false)
          item.autosaved = true
          setTimeout(() => { item.autosaved = false }, 2000)
        }
      }
    },

    async periodicalUpdate() {
      if (!this.activelist) return
      try {
        const result = await this.apiRequest('GET', `/api/refresh/${this.activelist.ID}/${this.lastUpdate}`)
        if (result === 'update') {
          this.selectList(this.activelist)
        }
      } catch (err) {
        // Ignore refresh errors
      }
    },

    copyItem(item) {
      navigator.clipboard.writeText(item.Title)
    },

    async changeListName() {
      if (!this.newListName) return
      await this.apiRequest('PUT', '/api/rename', { ID: this.activelist.ID, Name: this.newListName })
      this.activelist.Name = this.newListName
      this.newListName = ''
      this.renameListDialogVisible = false
    },

    async onSortEnd() {
      // vuedraggable already updated the array via v-model
      // Just save the new order to backend
      const sortOrder = this.todoList.map((item, index) => ({ id: item.ID, order: index }))
      await this.apiRequest('POST', '/api/sort', sortOrder)
    },

    async selectList(list) {
      if (!list) return
      this.activelist = list
      localStorage.setItem('activelist', JSON.stringify(list))

      try {
        const data = await this.apiRequest('GET', `/api/items/${list.ID}`)
        if (this.activelist.ID === this.account.Favorite) {
          this.activelist.Favorite = true
        }
        this.lastUpdate = Math.round(Date.now() / 1000)

        data.sort((a, b) => a.Order - b.Order)

        this.todoList = []
        this.completedToDoList = []
        const completed = []

        for (const item of data) {
          item.showNote = false
          item.showTime = false
          item.changed = false
          item.editing = false
          item.autosaved = false
          item.itemSize = 40
          item.dateTime = item.Time ? new Date(item.Time) : null
          item.RecurDays = item.RecurDays || 0
          item.PreAlarmMinutes = item.PreAlarmMinutes || 0
          item.TelegramChat = item.TelegramChat || ''

          if (!item.Complete) {
            this.todoList.push(item)
          } else {
            completed.push(item)
          }
        }

        completed.sort((a, b) => b.Completed - a.Completed)
        this.completedToDoList = completed
      } catch (err) {
        console.error(err)
      }
    },

    changeTitle(item) {
      item.editing = true
      item.originalTitle = item.Title
    },

    saveTitle(item) {
      item.editing = false
      if (item.Title !== item.originalTitle) {
        this.saveItem(item, false)
      }
    },

    cancelEditTitle(item) {
      item.Title = item.originalTitle
      item.editing = false
    },

    async fetchLists() {
      try {
        const data = await this.apiRequest('GET', '/api/lists')
        this.lists = data
        this.nolists = this.lists.length === 0

        if (this.lists.length > 0 && !this.activelist) {
          if (this.account.Favorite) {
            const favList = this.lists.find(l => l.ID === this.account.Favorite)
            if (favList) {
              favList.Favorite = true
              this.selectList(favList)
              return
            }
          }
          this.selectList(this.lists[0])
        } else if (this.activelist) {
          this.selectList(this.activelist)
        }
      } catch (err) {
        console.error(err)
      }
    },

    async addItem() {
      if (!this.entry) return
      const newEntry = {
        Note: '',
        Time: '',
        ListId: this.activelist.ID,
        Title: this.entry,
        Completed: 0,
        Complete: false,
        RecurDays: 0,
        PreAlarmMinutes: 0,
        TelegramChat: ''
      }

      try {
        const dbItem = await this.apiRequest('POST', '/api/item', newEntry)
        dbItem.showNote = false
        dbItem.showTime = false
        dbItem.changed = false
        dbItem.editing = false
        dbItem.autosaved = false
        dbItem.itemSize = 40
        dbItem.dateTime = null
        dbItem.RecurDays = 0
        dbItem.PreAlarmMinutes = 0
        dbItem.TelegramChat = ''
        this.todoList.unshift(dbItem)
        // Save order after adding new item at top
        const sortOrder = this.todoList.map((item, index) => ({ id: item.ID, order: index }))
        await this.apiRequest('POST', '/api/sort', sortOrder)
      } catch (err) {
        console.error(err)
      }
      this.entry = ''
    },

    showNote(item) {
      item.showNote = !item.showNote
      item.showTime = false
      item.itemSize = item.showNote ? 250 : 40
    },

    async saveItem(item, reload) {
      item.changed = false
      if (item.dateTime) {
        item.Time = item.dateTime.toISOString()
      }
      await this.apiRequest('PUT', '/api/item', item)
      if (reload) {
        this.selectList(this.activelist)
      }
    },

    saveNote(item) {
      this.saveItem(item, false)
      this.showNote(item)
    },

    saveTime(item) {
      this.saveItem(item, false)
      this.showTime(item)
    },

    clearTime(item) {
      item.Time = ''
      item.dateTime = null
      item.RecurDays = 0
      item.PreAlarmMinutes = 0
      item.TelegramChat = ''
      this.saveItem(item, false)
      this.showTime(item)
    },

    showTime(item) {
      item.showTime = !item.showTime
      item.showNote = false
      item.itemSize = item.showTime ? 280 : 40
    },

    completeItem(item) {
      item.Complete = true
      this.saveItem(item, true)
    },

    uncompleteItem(item) {
      item.Complete = false
      this.saveItem(item, true)
    }
  }
}
</script>

<style lang="scss" scoped>
.container {
  padding: 10px;
  min-height: calc(100vh - 20px);
}

.error {
  color: #fff;
  background: #c33;
  padding: 10px;
  border-radius: 5px;
}

.admin-text {
  color: #f55;
  font-weight: bold;
  margin-left: 10px;
}

.auth-form {
  max-width: 400px;
  margin: 20px auto;

  .form-field {
    display: flex;
    align-items: center;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 5px;
    margin-bottom: 5px;
    padding: 10px;

    i {
      color: #fff;
      font-size: 1.5rem;
      margin-right: 15px;
      width: 30px;
      text-align: center;
    }

    :deep(.p-inputtext) {
      flex: 1;
      background: transparent;
      border: none;
      color: #fff;
      font-size: 1.2rem;

      &::placeholder {
        color: rgba(255, 255, 255, 0.7);
      }
    }
  }
}

.menu-btn {
  position: fixed;
  top: 10px;
  right: 10px;
  z-index: 100;
}

.sidebar-menu {
  :deep(.p-drawer-content) {
    background: linear-gradient(to bottom right, #21656d, #0d4826);
    color: #fff;
  }

  :deep(.p-drawer-header) {
    background: rgba(0, 0, 0, 0.2);
    color: #fff;
  }
}

.list-menu-item {
  padding: 10px;
  cursor: pointer;
  font-size: 1.2rem;

  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .fav-icon {
    color: gold;
    margin-left: 5px;
  }
}

.sidebar-footer {
  position: absolute;
  bottom: 20px;
  left: 20px;
  right: 20px;

  .logged-in-text {
    display: block;
    font-size: 0.8rem;
    margin-bottom: 10px;
    opacity: 0.8;
  }
}

.list-header {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  background: rgba(0, 0, 0, 0.3);
  padding: 4px 8px;
  border-radius: 5px 5px 0 0;
  margin-bottom: 0;

  .list-name {
    color: #fff;
    font-size: 0.85rem;
    cursor: pointer;
    padding: 2px 5px;

    &:hover {
      text-decoration: underline;
    }
  }

  .action-btn {
    color: #fff;
    font-size: 0.75rem;
    cursor: pointer;
    padding: 3px 5px;
    border-radius: 3px;

    &:hover {
      background: rgba(255, 255, 255, 0.2);
    }

    &.delete-btn:hover {
      background: #a55;
    }

    &.share-btn:hover {
      background: #5a5;
    }

    &.fav-active {
      color: gold;
    }
  }
}

.add-item {
  display: flex;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 0 5px 5px 5px;
  overflow: hidden;

  .add-icon {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 60px;
    color: #fff;
    font-size: 1.5rem;
    cursor: pointer;
  }

  .add-input {
    flex: 1;
    background: transparent;
    border: none;
    color: #fff;
    font-size: 1.3rem;
    padding: 15px;

    &::placeholder {
      color: rgba(255, 255, 255, 0.7);
    }
  }

  :deep(.p-inputtext) {
    background: transparent;
    border: none;
    color: #fff;
    font-size: 1.3rem;

    &::placeholder {
      color: rgba(255, 255, 255, 0.7);
    }
  }
}

.todo-list {
  margin-top: 10px;

  .item {
    background: rgba(255, 255, 255, 0.9);
    border-radius: 5px;
    margin-bottom: 2px;
    overflow: hidden;
    transition: min-height 0.2s;

    &.expanded {
      overflow: visible;
    }
  }

  .item-content {
    display: grid;
    grid-template-columns: 40px 1fr auto auto;
    align-items: center;
    min-height: 40px;
    padding-right: 10px;
  }

  .item-checkbox {
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 1.2rem;
    color: #99c191;
    cursor: pointer;
  }

  .item-title {
    font-size: 1rem;
    padding: 5px;
    overflow: hidden;
    text-overflow: ellipsis;

    &.completed {
      text-decoration: line-through;
      opacity: 0.6;
    }

    .title-input {
      width: 100%;
      font-size: 1rem;
    }
  }

  .item-owner {
    font-size: 0.7rem;
    color: #888;
    margin-right: 10px;
  }

  .item-actions {
    display: flex;
    gap: 8px;

    .action-icon {
      cursor: pointer;
      color: #666;
      font-size: 1.1rem;

      &:hover {
        color: #333;
      }

      &.active {
        color: #c23741;
      }
    }
  }

  .item-panel {
    padding: 10px;
    background: rgba(255, 255, 255, 0.95);
    border-top: 1px solid #ddd;
  }

  .panel-actions {
    display: flex;
    gap: 5px;
    margin-top: 10px;
  }
}

.completed-list {
  .item {
    background: rgba(255, 255, 255, 0.5);
    opacity: 0.7;

    .item-title {
      text-decoration: line-through;
      color: #666;
    }

    .item-checkbox {
      color: #7a7;
    }
  }

  .item-panel {
    opacity: 1;
  }
}

.show-completed-wrapper {
  margin-top: 15px;
  text-align: center;
}

.show-completed-btn {
  display: inline-block;
  background: rgba(0, 0, 0, 0.4);
  color: #fff;
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.85rem;

  i {
    margin-right: 6px;
  }

  &:hover {
    background: rgba(0, 0, 0, 0.5);
  }
}

.clear-completed-wrapper {
  margin-top: 10px;
  text-align: center;
}

.drag-ghost {
  opacity: 0.5;
  background: rgba(255, 255, 255, 0.8) !important;
}

.drag-active {
  background: rgba(255, 255, 255, 0.95) !important;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
  border-radius: 5px;
  cursor: grabbing;
}

.alarm-options {
  margin: 10px 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.alarm-option {
  display: flex;
  align-items: center;
  gap: 8px;

  .alarm-label {
    color: #333;
    font-size: 0.85rem;
    min-width: 85px;
    font-weight: 500;
  }

  .alarm-input {
    width: 60px;

    :deep(.p-inputnumber-input) {
      width: 60px;
      padding: 4px 8px;
      text-align: center;
    }
  }

  .alarm-unit {
    color: #666;
    font-size: 0.8rem;
  }

  .telegram-select {
    min-width: 120px;

    :deep(.p-select-label) {
      padding: 4px 8px;
      font-size: 0.85rem;
    }
  }
}

.autosave-indicator {
  color: #5a5;
  font-size: 0.8rem;
  margin-left: auto;
  opacity: 0;
  transition: opacity 0.3s;

  &.visible {
    opacity: 1;
  }
}

.shared-users {
  margin-top: 20px;
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.drag-helper {
  background: rgba(255, 255, 255, 0.9) !important;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
  border-radius: 5px;
}

.mt-2 { margin-top: 0.5rem; }
.mt-3 { margin-top: 1rem; }
.mt-4 { margin-top: 1.5rem; }
.mb-3 { margin-bottom: 1rem; }
.mr-1 { margin-right: 0.25rem; }
.w-full { width: 100%; }
.text-center { text-align: center; }
</style>
