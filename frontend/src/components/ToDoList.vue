<template>
    <div class="container">
        <div v-if="errorText != ''">
            <span class="error">Something went wrong: {{ errorText }}
                <br>
                Please reload page.
            </span>
        </div>
        <div v-if="errorText == ''">
            <div v-if="account.Admin">
                <div>
                    <md-button class="md-default md-raised" id="adm-logout-btn" @click="logout()">Logout</md-button>
                    <span class="adminText">You are admin.</span>
                    <md-table v-model="accounts" md-sort="user" md-sort-order="asc" md-card>
                        <md-table-toolbar>
                            <h1 class="md-title">Accounts</h1>
                        </md-table-toolbar>
                        <md-table-row slot="md-table-row" slot-scope="{ item }">
                            <md-table-cell md-label="ID" md-numeric>{{ item.ID }}</md-table-cell>
                            <md-table-cell md-label="user" md-sort-by="name">{{ item.User }}</md-table-cell>
                            <md-table-cell md-label="email" md-sort-by="email">{{ item.Email }}</md-table-cell>
                            <md-table-cell md-label="admin" md-sort-by="admin">{{ item.Admin }}</md-table-cell>
                            <md-table-cell md-label="validated" md-sort-by="validated">{{ item.Validated }}</md-table-cell>
                            <md-table-cell md-label="updated" md-sort-by="updated">{{ item.CreatedAt }}</md-table-cell>
                            <md-table-cell md-label="created" md-sort-by="created">{{ item.UpdatedAt }}</md-table-cell>
                            <md-table-cell v-if="!item.Validated" md-label="validate"><md-button id="validate-btn" class="md-default md-raised" @click="validateAccount(item.ID)">Validate</md-button></md-table-cell>
                            <md-table-cell v-if="!item.Admin" md-label="remove"><md-button id="remove-btn" class="md-default md-raised" @click="removeAccount(item.ID)">Remove</md-button></md-table-cell>
                        </md-table-row>
                    </md-table>
                </div>
                <div>
                    <md-table v-model="tokens" md-sort="created" md-sort-order="asc" md-card>
                        <md-table-toolbar>
                            <h1 class="md-title">Tokens</h1>
                        </md-table-toolbar>
                        <md-table-row slot="md-table-row" slot-scope="{ item }">
                            <md-table-cell md-label="ID" md-numeric>{{ item.ID }}</md-table-cell>
                            <md-table-cell md-label="Account" md-sort-by="account">{{ item.Account.User }}</md-table-cell>
                            <md-table-cell md-label="Token" md-sort-by="token">{{ item.Token }}</md-table-cell>
                            <md-table-cell md-label="created" md-sort-by="created">{{ item.CreatedAt }}</md-table-cell>
                            <md-table-cell md-label="revoke"><md-button id="remove-btn" class="md-default md-raised" @click="revokeToken(item.ID)">Revoke</md-button></md-table-cell>
                        </md-table-row>
                    </md-table>
                </div>
            </div>
            <div id="signupdiv" v-if="createadm || createNewAcc">
                <span v-if="createadm" class="adminText"> This will create the admin account!</span>
                <div id="user" class="user-login">
                    <div class="user-icon">
                        <i class="fas fa-user"></i>
                    </div>
                    <form>
                        <input type="text" placeholder="Username" v-model="username">
                    </form>
                </div>
                <div id="pass" class="user-login">
                    <div class="user-icon">
                        <i class="fas fa-envelope"></i>
                    </div>
                    <form>
                        <input type="text" placeholder="email@example.com" v-model="email">
                    </form>
                </div>
                <div id="pass" class="user-login">
                    <div class="pass-icon">
                        <i class="fas fa-key"></i>
                    </div>
                    <form>
                        <input type="password" placeholder="Password" v-model="password">
                    </form>
                </div>
                <div id="pass" class="user-login">
                    <div class="pass-icon">
                        <i class="fas fa-key"></i>
                    </div>
                    <form>
                        <input placeholder="Repeat password..." v-on:keyup.enter="signup()" type="password" v-model="password2">
                    </form>
                </div>
                <md-button class="md-primary md-raised login-btn" @click="signup()">Create Account</md-button>
            </div>
            <div id="logindiv" v-if="!loggedIn && !createadm && !createNewAcc">
                <div id="user" class="user-login">
                    <div class="user-icon">
                        <i class="fas fa-user"></i>
                    </div>
                    <form>
                        <input v-on:keyup.enter="login()" type="text" placeholder="User" v-model="username">
                    </form>
                </div>
                <div id="pass" class="user-login">
                    <div class="pass-icon">
                        <i class="fas fa-key"></i>
                    </div>
                    <form @submit.prevent>
                        <input v-on:keyup.enter="login()" type="password" placeholder="Password" v-model="password">
                    </form>
                </div>
                <md-button class="md-primary md-raised login-btn" @click="login()">Login</md-button>
                <md-button id="signup-btn" class="md-primary signup-btn" @click="createNewAcc = true">Sign Up</md-button>
            </div>
            <div v-if="loggedIn && !account.Admin">
                <Slide :closeOnNavigation="true" >
                <md-dialog-prompt
                    :md-active.sync="newListActive"
                    v-model="newListValue"
                    md-title="New List Name"
                    md-inp2t-maxlength="20"
                    @md-confirm="newList"
                    md-confirm-text="Create"/>
                    <md-button class="md-primary md-raised" @click="newListActive = true">New List</md-button>
                    <div v-for="item in lists" :key="item.ID">
                        <a href="#"> 
                            <span class="list-menu-item" @click="selectList(item)">{{ item.Name }}<span v-if="item.ID == account.Favorite" class="fa fa-star fav-size"> </span></span>
                        </a>
                    </div>
                    <span style="font-size: 0.8rem;">Logged in as <span style="font-size: 0.8rem; color: #fff;">{{ account.User }}</span></span>
                    <md-button class="md-default md-raised" @click="logout()">Logout</md-button>
                </Slide>
                <div v-if="nolists"> 
                    <md-button class="md-primary md-raised" @click="newListActive = true">Create your first list</md-button>
                </div>
                <div id="listid" class="list-name" v-if="activelist">
                    <div class="button">
                        <md-dialog-prompt
                            :md-active.sync="listActive"
                            v-model="newListName"
                            md-title="Change List Name"
                            md-input-maxlength="20"
                            @md-confirm="changeListName"
                            md-confirm-text="Update"/>

                            <md-button id="changeListName" class="md-primary" @click="listActive = true">{{ activelist.Name }} </md-button>
                            <md-button class="md-primary" id="delete-btn" @click="showDeleteList = true"><span class="fa fa-trash"></span></md-button>
                            <md-button class="md-primary" id="share-btn" @click="showShareList = true"><span class="fa fa-user-plus"></span></md-button>
                            <md-button @click="setFavorite()" id="fav-btn" class="md-primary" :class="{'fav-btn-active': activelist.Favorite == true}"><span class="fa fa-star"></span></md-button>
                            <md-button id="addToHome" class="md-primary md-raised"><span class="fa fa-mobile-alt"></span></md-button>


                            <md-dialog-confirm
                                :md-active.sync="showDeleteList"
                                md-title="Delete List?"
                                md-content="This will remove the list AND all items in the list, are you really sure?"
                                @md-confirm="deleteList()"
                                md-confirm-text="Delete"/>

                                <md-dialog-prompt
                                    :md-active.sync="showShareList"
                                    v-model="shareWithUser"
                                    md-title="Share List With User"
                                    md-content="Enter the e-mail or user you want to share this list with?"
                                    md-inp2t-maxlength="20"
                                    @md-confirm="shareList()"
                                    md-confirm-text="Share"/>
                    </div>
                </div>
                <div class="add-item" v-if="!errorText && activelist">
                    <div class="add-icon" @click="addItem">
                        <i class="fas fa-plus"></i>
                    </div>
                    <form @submit.prevent="addItem">
                        <input type="text" placeholder="Add..." v-model="entry">
                    </form>
                </div>
                <div>
                </div>
                <div class="todo-list list-list" id="todolist">
                    <SlickList :shouldCancelStart="cancelSort" @sort-start="sortStart($event)" helperClass="drag-helper" :transitionDuration="0" :lockToContainerEdges="true" :pressDelay="350" @input="sortEnd" lockAxis="y" v-model="todoList">
                    <SlickItem v-on:dblclick.native="changeTitle(item)" :style="{height: item.itemSize+'px'}" :index="index" class="item" :class="{'show': item.show}" v-for="(item, index) in todoList" :key="item.ID">
                    <div class="item-checkbox">
                        <i v-if="!item.Complete" class="fas fa-square" @click="completeItem(item)"></i>
                        <i v-else class="fas fa-check-square"></i>
                    </div>
                    <div class="item-title" :id="item.ID">
                        {{ item.Title }} 
                    </div>
                    <div class="item-owner" v-if="activelist.Share != null && activelist.Share.length > 0">
                        <span style="margin-right: 10px; font-size: 0.7rem; color: #888;"> [{{item.Account.User}}]</span>
                    </div>
                    <div class="feature-icon-copy">
                        <i class="fas fa-copy" @click="copy(item);"></i>
                    </div>
                    <div class="feature-icon" :class="{'red-bg': item.Note.length > 0}">
                        <i class="fas fa-pen" @click="showNote(item);"></i>
                    </div>
                    <div class="feature-icon-bell" :class="{'red-bg': item.Time.length > 0}">
                        <i class="fas fa-bell" @click="showTime(item)"></i>
                    </div>

                    <div class="item-note" v-if="item.showNote">
                        <md-field>
                            <label>Note</label>
                            <md-textarea v-model="item.Note" @input="noteInput(item)"></md-textarea>
                        </md-field>
                        <md-button class="md-primary md-raised" @click="saveNote(item)">Save</md-button>
                        <md-button class="md-accent md-raised" @click="showNote(item)">Cancel</md-button>
                    </div>
                    <div class="item-alarm" v-if="item.showTime">
                        <datetime class="vdatetime-input" type="datetime" v-model="item.Time"></datetime>
                        <br>
                        <md-button class="md-primary md-raised" @click="saveTime(item)">Set Alarm</md-button>
                        <md-button class="md-accent md-raised" @click="clearTime(item);">Clear</md-button>
                    </div>
                    </SlickItem>
                    </SlickList>
                </div>
                <div class="show-completed" v-if="completedToDoList.length > 0">
                    <div class="button" @click="showCompletedList = !showCompletedList">
                        <span v-if="!showCompletedList">show</span><span v-else>hide</span> 
                        completed ({{ completedToDoList.length }})
                    </div>
                </div>
                <div class="todo-list complete-list" v-if="showCompletedList">
                    <div :style="{height: item.itemSize+'px'}" class="item" :class="{'show': item.show}" v-for="item in completedToDoList" :key="item.ID">
                        <div class="item-checkbox">
                            <i v-if="!item.Complete" class="fas fa-square"></i>
                            <i v-else class="fas fa-check-square" @click="uncompleteItem(item)"></i>
                        </div>
                        <div class="item-title">
                            {{ item.Title }}
                        </div>
                        <div class="item-owner" v-if="activelist.Share != null && activelist.Share.length > 0">
                            <span style="margin-right: 10px; font-size: 0.7rem; color: #888;"> [{{item.Account.User}}]</span>
                        </div>
                        <div class="feature-icon-copy">
                            <i class="fas fa-copy" @click="copy(item);"></i>
                        </div>
                        <div class="feature-icon" :class="{'red-bg': item.Note.length > 0}">
                            <i class="fas fa-pen" @click="showNote(item);"></i>
                        </div>
                        <div class="feature-icon-bell" :class="{'red-bg': item.Time.length > 0}">
                            <i class="fas fa-bell" @click="showTime(item)"></i>
                        </div>
                        <div class="item-note" v-if="item.showNote">
                            <md-field>
                                <label>Note</label>
                                <md-textarea v-model="item.Note" @input="noteInput(item)"></md-textarea>
                            </md-field>
                            <md-button class="md-primary md-raised" @click="saveNote(item)">Save</md-button>
                            <md-button class="md-accent md-raised" @click="showNote(item)">Cancel</md-button>
                        </div>
                        <div class="item-alarm" v-if="item.showTime">
                            <datetime class="vdatetime-input" type="datetime" v-model="item.Time"></datetime>
                            <br>
                            <md-button class="md-primary md-raised" @click="saveTime(item)">Set Alarm</md-button>
                            <md-button class="md-accent md-raised" @click="clearTime(item);">Clear</md-button>
                        </div>
                    </div>
                    <md-button class="md-accent md-raised" @click="clearCompleted();">Clear Completed</md-button>
                </div>
                <div v-if="activelist != undefined && activelist != ''">
                    <md-dialog-confirm
                        :md-active.sync="showRemoveShareList"
                        v-model="removeUser"
                        md-title="Remove Share?"
                        md-content="This will the user from shared list, are you really sure?"
                        @md-confirm="removeUserShare()"
                        md-confirm-text="Remove"/>
                        <div style="margin-top: 20px;" v-if="activelist.Share != null">
                            <div v-if="activelist.Share.length > 0 ">
                                <md-button class="md-primary" id="owner-btn"><span class="fa fa-user"> {{activelist.Account.User}}</span></md-button>
                                <span :index="index" v-for="(item, index) in activelist.Share" :key="item.ID">
                                    <md-button @click="removeUser=item; showRemoveShareList = true" class="md-primary" id="shared-btn"><span class="fa fa-users"> {{item.User}} </span></md-button>
                                </span>
                            </div>
                        </div>
                </div>
            </div>
        </div>
    </div>
</template>


<script>
import { SlickList, SlickItem } from 'vue-slicksort';
import { Slide } from 'vue-burger-menu'
import $ from 'jquery';

//import axios from 'axios';
export default {
    name: 'app',
    data: () => {
        return {
            account: "",
            tokens: [],
            accounts: [],
            nolists: false,
            createNewAcc: false,
            email: "",
            createadm: false,
            username: "",
            password: "",
            password2: "",
            removeUser: "",
            shareWithUser: "",
            showRemoveShareList: false,
            showShareList: false,
            showDeleteList: false,
            newListActive: false,
            newListValue: null,
            loggedIn: false,
            lastUpdate: '',
            timer: '',
            newListName: "",
            listActive: false,
            activelist: "",
            itemSize: 40,
            result: "",
            errorText: "",
            info: "",
            entry: '',
            entryAlarm: false,
            entryEdit: false,
            lists: [],
            todoList: [],
            completedToDoList: [],
            isSorting: false,
            showCompletedList: false
        }
    },
    components: {
        SlickItem,
        SlickList,
        Slide
    },
    created: function() {
        // Check if already logged in, not that safe, but...
        var that = this;
        if (localStorage.token != "" && localStorage.token != undefined) {
            this.loggedIn = true;
            this.account = JSON.parse(localStorage.getItem('account'));
            this.selectList(JSON.parse(localStorage.getItem('activelist')));

            if(this.account.Admin == true) {
                this.fetchAdminItems();
            } else {
                this.fetchLists()
                setInterval(function() {
                    if(that.loggedIn) {
                        that.periodcalUpdate();
                    }
                }, 5000);
            }
        } else {
            // Check if admin has been created.
            $.ajax({
                   type: "GET",
                   url: "/hasadm",
                   success: function(data) {
                       if(data == false) {
                           that.createadm = true;
                       }else {
                           that.createadm = false;
                       }
                   },
                   error: that.handleError
            });
        }
    },
    methods: {
        // Remove a token for a user, admin operation
        revokeToken(id) {
            var that = this;
            $.ajax({
                   type: "PUT",
                   url: "/api/removetoken/"+id,
                   beforeSend: that.setHeader,
                   success: function() {
                       that.fetchAdminItems();
                   },
                   error: that.handleError
            });
        },
        // Validate a signed up user, admin operation
        validateAccount(id) {
            var that = this;
            $.ajax({
                   type: "PUT",
                   url: "/api/validate/"+id,
                   beforeSend: that.setHeader,
                   success: function() {
                       that.fetchAdminItems();
                   },
                   error: that.handleError
            });
        },
        // Remove account, admin operation
        removeAccount(id) {
            var that = this;
            $.ajax({
                   type: "PUT",
                   url: "/api/removeaccount/"+id,
                   beforeSend: that.setHeader,
                   success: function() {
                       that.fetchAdminItems();
                   },
                   error: that.handleError
            });
        },
        // Get all accounts and tokens, admin operation
        fetchAdminItems() {
            var that = this;
            this.tokens = [];
            this.accounts = [];
            $.ajax({
                   type: "GET",
                   url: "/api/allaccounts",
                   beforeSend: that.setHeader,
                   success: function(data) {
                       that.accounts = data;
                       $('#app').css("max-width", "100%");
                       $('#app').css("width", "100%");
                   },
                   error: that.handleError
            });

            $.ajax({
                   type: "GET",
                   url: "/api/alltokens",
                   beforeSend: that.setHeader,
                   success: function(data) {
                       that.tokens = data;
                   },
                   error: that.handleError
            });

        },
        // Signup for a new account
        signup() {
            if(this.email != "" && this.password2 != "" && this.password != "" && this.username != "" && (this.password == this.password2)) {
                var data = {};
                data["user"] = this.username;
                data["pass"] = this.password;
                data["adm"] = this.createadm;
                data["email"] = this.email;
                var that = this;
                this.email = "";
                this.password = "";
                this.password2 = "";
                this.createadm = false;

                $('#signupdiv').html("<img style='padding-left: 50%;' src='loader.gif'/>");
                $.ajax({
                       type: "POST",
                       url: "/signup",
                       data: JSON.stringify(data),
                       success: function() {
                           window.location.reload();
                       },
                       error: that.handleError
                });
            }
        },
        // Set the current list as favorite
        setFavorite() {
            this.activelist.Favorite = !this.activelist.Favorite;
            var fav = {}
            fav["list"] = this.activelist.ID;
            fav["favorite"] = this.activelist.Favorite;

            if (this.activelist.Favorite) {
                this.account.Favorite = this.activelist.ID;
            } else {
                this.account.Favorite = 0;
            }

            localStorage.setItem('account', JSON.stringify(this.account));

            var that = this;
            $.ajax({
                   type: "POST",
                   url: "/api/favorite",
                   data: JSON.stringify(fav),
                   beforeSend: that.setHeader,
                   success: function() {
                       that.fetchLists()
                   },
                   error: that.handleError
            });
        },
        // Remove a user from a shared list
        removeUserShare() {
            if(this.removeUser == "") {
                return;
            }
            var user = this.removeUser;
            this.removeUser = "";

            var that = this;
            var share = {};
            share["user"] = user.User;
            share["list"] = this.activelist.ID;

            $.ajax({
                   type: "POST",
                   url: "/api/removesharelist",
                   data: JSON.stringify(share),
                   beforeSend: that.setHeader,
                   success: function() {
                       that.fetchLists()
                       let index = that.activelist.Share.findIndex(element => element.ID === user.ID);
                       that.activelist.Share.splice(index, 1);
                   },
                   error: that.handleError
            });
        },
        // Add a user to share a list with
        shareList() {
            var that = this;
            var share = {};
            var user = this.shareWithUser;
            share["user"] = user;
            share["list"] = this.activelist.ID;

            $.ajax({
                   type: "POST",
                   url: "/api/sharelist",
                   data: JSON.stringify(share),
                   beforeSend: that.setHeader,
                   success: function(userData) {
                       that.fetchLists()
                       if (that.activelist.Share == null) {
                           that.activelist.Share = [];
                       }
                       if (userData != "") {
                           that.activelist.Share.push(userData);
                       }
                       that.shareWithUser = "";
                   },
                   error: that.handleError
            });
        },
        // Clear completed items for list 
        clearCompleted() {
            var that = this;
            $.ajax({
                   type: "POST",
                   url: "/api/deletecompleted",
                   data: JSON.stringify(this.activelist),
                   beforeSend: that.setHeader,
                   success: function() {
                       that.selectList(that.activelist);
                       that.showCompletedList = false;
                   },
                   error: that.handleError
            });
        },
        // Delete the current list
        deleteList() {
            var that = this;
            if(this.account.Favorite == this.activelist.ID) {
                this.account.Favorite = 0;
            }
            $.ajax({
                   type: "POST",
                   url: "/api/deletelist",
                   data: JSON.stringify(this.activelist),
                   beforeSend: that.setHeader,
                   success: function() {
                       that.todoList = [];
                       that.completedToDoList = [];
                       that.activelist = "";
                       that.fetchLists();

                   },
                   error: that.handleError
            });
        },
        // Log out the user
        logout() {
            var that = this;
            $.ajax({
                   type: "POST",
                   url: "/api/logout",
                   beforeSend: that.setHeader,
            });
            that.loggedIn = false;
            that.account = "";
            localStorage.clear();
            $('#app').css("max-width", "768px");
            $('#app').css("width", "100%");
        },
        // Handle any error from backend
        handleError(err) {
            console.log(err);
            //this.errorText = err;
            //if(err.status == 401) {
            //    this.loggedIn = false;
            //    localStorage.clear();
            //    window.location.reload();
            //}
            window.location.reload();
        },
        // Set the token in the header for a request
        setHeader(request) {
            request.setRequestHeader("token", localStorage.token);
        },
        // Log in the user
        login() {
            if(this.password != "" && this.username != "") {
                var data = {};
                data["user"] = this.username;
                data["pass"] = this.password;
                var that = this;
                this.password = "";

                $('#logindiv').html("<img style='padding-left: 50%;' src='loader.gif'/>");
                $.ajax({
                       type: "POST",
                       url: "/login",
                       data: JSON.stringify(data),
                       success: function(data) {
                           that.loggedIn = true;
                           localStorage.token = data.Token;
                           localStorage.setItem('account', JSON.stringify(data.Account));
                           that.account = data.Account;
                           if (that.account.Admin == true) {
                               that.fetchAdminItems();
                           } else {
                               that.fetchLists();
                           }
                       },
                       error: that.handleError
                });
            }
        },
        // Create a new list
        newList() {
            var newEntry = {
                Name: this.newListValue
            };
            this.newListValue = "";

            var that = this;
            if(newEntry.Name != "" && newEntry.Name != null) {
                $.ajax({
                       type: "POST",
                       beforeSend: that.setHeader,
                       url: "/api/list",
                       data: newEntry,
                       success: function(data) {
                           that.activelist = data;
                           that.fetchLists();
                       },
                       error: that.handleError
                });
            }
        },
        // Check if there have been any updates to current list
        periodcalUpdate() {
            if(this.activelist == undefined || this.activelist == "" || this.isSorting) {
                return;
            }
            var that = this;
            $.ajax({
                   type: "GET",
                   beforeSend: that.setHeader,
                   url: "/api/refresh/"+that.activelist.ID+"/"+that.lastUpdate,
                   success: function(data) {
                       if (data == "update") {
                           that.selectList(that.activelist);
                       }

                       if (that.errorText != "") {
                           that.errorText = "";
                           window.location.reload();
                       }
                   },
                   error: that.handleError
            });
        },
        // Copy is used to copy clicked items title
        copy(item) {
            var temp = $("<input>");
            $("body").append(temp);
            temp.val(item.Title);
            temp.select();
            document.execCommand("copy");
            temp.remove();
        },
        // Change name of the current list
        changeListName() {
            var list = {};
            list.ID = this.activelist.ID;
            if(this.newListName == "") {
                return;
            }
            list.Name = this.newListName;
            var that = this;
            $.ajax({
                   type: "PUT",
                   beforeSend: that.setHeader,
                   url: "/api/rename",
                   data: JSON.stringify(list),
                   success: function() {
                       that.activelist.Name = list.Name;
                   },
                   error: that.handleError
            });
            this.newListName = "";
        },
        // Dont allow dragging while having date selector up
        cancelSort() {
            if ($(".vdatetime-popup").length != 0) {
                return true;
            }
            return false;
        },
        // Sorting event for moving an item in the list
        sortStart() {
            this.isSorting = true;
            let newCanvas = document.querySelector('.drag-helper')
            var text = newCanvas.innerText;
            newCanvas.innerHTML = "<div style="+
                "'font-size: 1.0rem;"+
                "background: hsla(0,0%,100%,.5);"+
                "border-radius: 5px;"+
                "padding-top: 10px;"+
                "padding-left: 40px;"+
                "width: 100%;"+
                "height: 40px;"+
                "margin-bottom: 2px;"+
                "overflow:hidden;'"+
                ">"+text+"</div>";
        },
        // Sorting event end for moving an item in the list
        sortEnd(a) {
            this.isSorting = false;
            var sortOrder = [];
            for(var i = 0; i < a.length; i++) {
                sortOrder.push({"id": a[i].ID, "order": i})
            }
            var that = this;
            $.ajax({
                   type: "POST",
                   beforeSend: that.setHeader,
                   url: "/api/sort",
                   data: JSON.stringify(sortOrder),
                   success: function() {
                       that.selectList(that.activelist);
                   },
                   error: that.handleError
            });
        },
        // Select a list and fetch its items
        selectList(list) {
            if(list == undefined || list == "") {
                return;
            }
            this.activelist = list;
            var that = this;
            localStorage.setItem('activelist', JSON.stringify(list));

            $.ajax({
                   type: "GET",
                   beforeSend: that.setHeader,
                   url: "/api/items/"+list.ID,
                   success: function(data) {
                       if (that.activelist.ID == that.account.Favorite) {
                           that.activelist.Favorite = true;
                       }
                       that.lastUpdate = Math.round(new Date().getTime()/1000);

                       data.sort(function(a, b){
                           return a.Order-b.Order;
                       })

                       that.todoList = [];
                       that.completedToDoList = [];
                       var completed = [];

                       that.activelist.Items = [];
                       for(var k in data) {
                           var v = data[k];
                           v.showNote = false;
                           v.showTime = false;
                           v.show = true;

                           if(!v.Complete) {
                               that.todoList.push(v);
                           } else {
                               completed.push(v); 
                           }
                           that.activelist.Items.push(v);
                       }
                       // Sort completed.
                       completed.sort(function(a,b){
                           return b.Completed - a.Completed;
                       });
                       that.completedToDoList = completed;
                   },
                   error: that.handleError
            });
        },
        // Change the title of an item
        changeTitle(item) {
            var that = this;
            $('#'+item.ID).html("<input id='inp_"+item.ID+"' style='font-size: 1.2rem; background: none; width: 100%; height: 50px; color: #000; box-shadow: none; border: none;' value='"+item.Title+"'></input>");
            $('#inp_'+item.ID).keyup(function(e){
                // Cancel
                if (e.keyCode == 27) {
                    $('#'+item.ID).text(item.Title);
                }
                // Save
                if(e.keyCode == 13) {
                    item.Title = $('#inp_'+item.ID).val();
                    $('#'+item.ID).text(item.Title);
                    that.saveItem(item);
                }
            });
        },
        // Fetch all lists
        fetchLists() {
            var that = this;
            $.ajax({
                   type: "GET",
                   beforeSend: that.setHeader,
                   url: "/api/lists",
                   success: function(data) {
                       that.updateLists(data);
                   },
                   error: that.handleError
            });
        },
        // Update list data
        updateLists(data) {
            this.lists = data;
            this.nolists = false;
            if (this.lists.length == 0) {
                this.nolists = true;
            }

            if (this.lists.length > 0 && (this.activelist == undefined || this.activelist == "")) {
                // Find favorite
                if (this.account.Favorite != 0) {
                    for(var i = 0; i < this.lists.length; i++) {
                        if (this.lists[i].ID == this.account.Favorite) {
                            this.lists[i].Favorite = true;
                            this.selectList(this.lists[i]);
                            return;
                        }
                    }
                    this.selectList(this.lists[0]);
                } else {
                    this.selectList(this.lists[0]);
                }
            } else {
                this.selectList(this.activelist);
            }
        },
        // Add a new entry
        addItem() {
            if(this.entry !== '') {
                let newEntry = {
                    Note: "",
                    Time: "",
                    ListId: "",
                    Title: this.entry,
                    Completed: 0,
                    Complete: false,
                    showNote: false,
                    showTime: false,
                    show: true,
                }

                newEntry.ListId = this.activelist.ID;
                var that = this;
                $.ajax({
                       type: "POST",
                       beforeSend: that.setHeader,
                       url: "/api/item",
                       data: newEntry,
                       success: function(dbItem) {
                           dbItem.showNote = false;
                           dbItem.showTime = false;
                           dbItem.show = true;
                           that.todoList.splice(0, 0, dbItem);
                           that.sortEnd(that.todoList);
                       },
                       error: that.handleError
                });
            }

            this.entry = '';
        },
        // Show the note for an item
        showNote(item) {
            item.showNote = !item.showNote;
            item.alarm = false;
            item.showTime = false;
            if (item.showNote) {
                item.itemSize = 400;
            } else {
                item.itemSize = 40;
            }
        },
        noteInput() {
            // TBD: Save each entry...? "auto saved.." timer?
        },
        // Save an item with any changes made
        saveItem(item) {
            var that = this;
            $.ajax({
                   type: "PUT",
                   beforeSend: that.setHeader,
                   url: "/api/item",
                   data: item,
                   success: function() {
                       that.selectList(that.activelist);
                   },
                   error: that.handleError
            });
        },
        // Save a note for an item
        saveNote(item) {
            this.saveItem(item);
            this.showNote(item);
        },
        // Save a alarm for an item
        saveTime(item) {
            this.saveItem(item);
            this.showTime(item);
        },
        // Clear the alarm for an item
        clearTime(item) {
            item.Time = "";
            this.saveItem(item);
            this.showTime(item);
        },
        // Show the alarm settings
        showTime(item) {
            item.showTime = !item.showTime;
            item.showNote = false;
            if (item.showTime) {
                item.itemSize = 200;
            } else {
                item.itemSize = 40;
            }
        },
        // Mark an items as complete
        completeItem(item) {
            item.Complete = !item.Complete;
            item.show = false;
            this.saveItem(item);
        },
        // Mark an item as not complete
        uncompleteItem(item) {
            item.Complete = !item.Complete;
            item.show = false;
            this.saveItem(item);
        } 
    }
}

</script>

<style lang="scss" scoped>
.container {
    padding: 10px;
    margin-left: 0px;
    width: calc(100% - 1px);
    min-height: calc(100% - 10px);

    @mixin feature-icon-copy($color) {
        grid-column-start: 4;
        grid-column-end: 4;
        display: flex;
        justify-content: center;
        align-items: center;
        color: #719fbb;
        font-size: 1.4rem;
        margin: 0 0px;
        opacity: 0.2;
        cursor: pointer;
    }

    @mixin feature-icon($color) {
        grid-column-start: 5;
        grid-column-end: 5;
        display: flex;
        justify-content: center;
        align-items: center;
        color: $color;
        font-size: 1.4rem;
        margin: 0 0px;
        cursor: pointer;
    }

    @mixin feature-icon-bell($color) {
        grid-column-start: 6;
        grid-column-end: 6;
        display: flex;
        justify-content: center;
        align-items: center;
        color: $color;
        font-size: 1.4rem;
        margin: 0 0px;
        cursor: pointer;
    }

    .error {
        color: #FFF;
        font-size: 0.7rem;
        background: #CC3333;
        margin-top: -10px;
        margin-bottom: -10px;
    }

    .add-item {
        display: grid;
        grid-template-columns: 70px auto;
        grid-template-rows: 60px;
        width: 100%;
        height: 60px;
        background: rgba($color: #000000, $alpha: .3);
        border-radius: 0px 5px 5px 5px;
        overflow: hidden;

        .add-icon {
            grid-column-start: 1;
            grid-column-end: 1;
            display: flex;
            justify-content: center;
            align-items: center;
            color: #fff;
            font-size: 2rem;
            cursor: pointer;
        }

        input {
            grid-column-start: 2;
            grid-column-end: 2;
            background: none;
            border: none;
            width: 100%;
            height: 60px;
            color: #fff;
            font-size: 1.6rem;

            &::placeholder {
                color: #fff;
            }

            &:focus {
                outline: none;
            }
        }

    }

    .drag{
        font-size: 25px;
        background: rgba($color: #fff, $alpha: .8);
        border-radius: 5px;
        color: #FF000;
        display: grid;
    }

    .todo-list {
        padding-top: 20px;

        .item {
            display: grid;
            grid-template-columns: 40px auto 30px 30px 30px 30px;
            grid-template-rows: 40px;
            width: 100%;
            height: 40px;
            margin-bottom: 2px;
            background: rgba($color: #fff, $alpha: .8);
            border-radius: 5px;
            overflow: hidden;
            opacity: 0;
            transition: all .2s linear;

            .item-alarm {
                grid-column-start: 1;
                grid-column-end: 5;
                width: 100%;
                height: 98%;
            }
            .item-note {
                grid-column-start: 1;
                grid-column-end: 7;
                width: 100%;
                height: 78%;
            }
            .md-field {
                height: 98%;
            }
            .md-field .md-textarea {
                min-height: 98%;
            }

            .item-checkbox {
                grid-column-start: 1;
                grid-column-end: 1;
                display: flex;
                justify-content: center;
                align-items: center;
                font-size: 1.2rem;
                color: #99C191;
                cursor: pointer;
            }

            .item-title {
                -webkit-user-select: none;
                -moz-user-select: none;
                -ms-user-select: none;
                user-select: none;
                grid-column-start: 2;
                grid-column-end: 2;
                display: flex;
                justify-content: flex-start;
                align-items: center;
                font-size: 1.0rem;
                overflow: hidden;
            }

            .item-owner {
                grid-column-start: 3;
                grid-column-end: 3;
                display: flex;
                justify-content: center;
                align-items: center;
                color: #719fbb;
                font-size: 0.8rem;
                margin: 0 0px;
            }

            .feature-icon-copy {
                @include feature-icon-copy(#333);
            }
            .feature-icon {
                @include feature-icon(#333);
            }
            .feature-icon-bell {
                @include feature-icon-bell(#333);
            }
        }

        .show {
            opacity: 1;
        }
    }

    .complete-list {
        .item {
            position: relative;
            background: rgba($color: #fff, $alpha: .4);

            &::before {
                content: '';
                position: absolute;
                top: 20px;
                width: calc(100% - 120px);
                margin: 0 40px;
                border-bottom: 2px solid #555;
            }
        }
    }

    .show-completed {
        display: flex;
        justify-content: center;
        align-items: center;
        padding-top: 20px;

        .button {
            color: #fff;
            padding: 10px;
            font-size: 1.0rem;
            text-transform: uppercase;
            background: rgba($color: #000000, $alpha: .3);
            border-radius: 5px;
            overflow: hidden;
            cursor: pointer;
        }
    }

    .datetime  {
        margin-left: 20px;
        margin-top: 5px;
        border-radius: 10px;
    }
    .red-bg {
        color: #C23741 !important;
    }
}

.time {
    font-size: 0.7rem;
    justify-content: right;
    color: #333333;
}

.list-name{
    display: inline-flex;
    justify-content: right;
    align-items: center;
    padding-top: 0px;
    font-weight: bold;

    .button {
        color: #fff;
        padding: 5px;
        font-size: 0.8rem;
        background: rgba($color: #000000, $alpha: .3);
        border-radius: 5px 5px 0px 0px;
        overflow: hidden;
        margin-top: 10px;
        z-index: 0;
    }
}
.vdatetime-input { 
    font-size: 1rem;
    margin: 10px;
}

#addToHome{
    height: 20px;
    margin: 1px;
    background: #13668f;
    color: #FFF;
    font-size: 10px;
    text-transform: none;
    display: none;
    min-width: 20px;
    width: 20px;
}
#changeListName {
    color: #fff;
    height: 20px;
    margin: 1px;
    font-size: 11px;
    text-transform: none;
}


.user-login {
    display: grid;
    grid-template-columns: 70px auto;
    grid-template-rows: 60px;
    width: 100%;
    height: 60px;
    background: rgba($color: #000000, $alpha: .3);
    border-radius: 5px 5px 5px 5px;
    overflow: hidden;

    .pass-icon {
        grid-column-start: 1;
        grid-column-end: 1;
        display: flex;
        justify-content: center;
        align-items: center;
        color: #fff;
        font-size: 2rem;
        cursor: pointer;
    }
    .user-icon {
        grid-column-start: 1;
        grid-column-end: 1;
        display: flex;
        justify-content: center;
        align-items: center;
        color: #fff;
        font-size: 2rem;
        cursor: pointer;
    }

    input {
        grid-column-start: 2;
        grid-column-end: 2;
        background: none;
        border: none;
        width: 100%;
        height: 60px;
        color: #fff;
        font-size: 1.6rem;

        &::placeholder {
            color: #fff;
        }

        &:focus {
            outline: none;
        }
    }
}
.login-btn {
    width: 100%;
    margin: 0px;
    padding: 0px;
}

#signup-btn {
    width: 100px;
    background: transparent !important;
    color: #FFF;
    padding: 0px;
    margin-top: 20px;
}

#actions {
    justify-content: center;
}

.fav-btn-active {
    background: #FFD700 !important;
}
#fav-btn {
    background: transparent;
    color: #FFFFFF;
    height: 20px;
    min-height: 20px;
    width: 20px;
    min-width: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
}

#delete-btn {
    background: #AA5555;
    color: #FFFFFF;
    height: 20px;
    min-height: 20px;
    width: 20px;
    min-width: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
}

#share-btn {
    background: #55AA55;
    color: #FFFFFF;
    height: 20px;
    min-height: 20px;
    width: 20px;
    min-width: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
    font-family: Arial,Helvetica,sans-serif;
}

#creator-btn {
    font-family: Arial,Helvetica,sans-serif;
    background: transparent;
    color: #000;
    height: 20px;
    min-height: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
}

#owner-btn {
    font-family: Arial,Helvetica,sans-serif;
    background: transparent;
    color: #FFFFFF;
    height: 20px;
    min-height: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
}

#adm-logout-btn {
    background: #AA5555;
    color: #FFFFFF;
    height: 20px;
    min-height: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
}

#remove-btn {
    background: #AA5555;
    color: #FFFFFF;
    height: 20px;
    min-height: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
}

#validate-btn {
    background: #55AA55;
    color: #FFFFFF;
    height: 20px;
    min-height: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
}

#shared-btn {
    background: #55AA55;
    color: #FFFFFF;
    height: 20px;
    min-height: 20px;
    margin: 1px;
    font-size: 10px;
    text-transform: none;
}

#user {
    border-radius: 5px 5px 0px 0px;
}
#pass {
    border-radius: 0px;
}

.fav-size {
    color: #FFD700;
    padding-left: 5px;
}

.adminText {
    color: #FF5555;
    text-decoration: bold;
}

</style>
