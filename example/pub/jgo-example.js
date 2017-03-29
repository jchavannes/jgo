JGoExampleApp = {};

$(function() {

    /**
     * @param token {string}
     */
    function InitCsrf(token) {
        /**
         * @param method {string}
         * @returns {boolean}
         */
        function csrfSafeMethod(method) {
            // HTTP methods that do not require CSRF protection.
            return (/^(GET|HEAD|OPTIONS|TRACE)$/.test(method));
        }

        $.ajaxSetup({
            crossDomain: false,
            beforeSend: function (xhr, settings) {
                if (!csrfSafeMethod(settings.type)) {
                    xhr.setRequestHeader("X-CSRF-Token", token);
                }
            }
        });
    }

    JGoExampleApp.InitCsrf = InitCsrf;

    var BaseURL = "/";

    /**
     * @param url {string}
     */
    function SetBaseUrl(url) {
        BaseURL = url;
    }

    JGoExampleApp.SetBaseUrl = SetBaseUrl;

    JGoExampleApp.Form = {
        /**
         * @param {jQuery} $ele
         */
        Signup: function($ele) {
            $ele.submit(function (e) {
                e.preventDefault();
                var username = $ele.find("[name=username]").val();
                var password = $ele.find("[name=password]").val();

                if (username.length == 0) {
                    alert("Must enter a username.");
                    return;
                }

                if (password.length == 0) {
                    alert("Must enter a password.");
                    return;
                }

                $.ajax({
                    type: "POST",
                    url: BaseURL + JGoExampleApp.URL.SignupSubmit,
                    data: {
                        username: username,
                        password: password
                    },
                    success: function () {
                        window.location = BaseURL + JGoExampleApp.URL.Lobby
                    },
                    /**
                     * @param {XMLHttpRequest} xhr
                     */
                    error: function (xhr) {
                        alert("Error creating account:\n" + xhr.responseText + "\nIf this problem persists, try refreshing the page.");
                    }
                });
            });
        },
        /**
         * @param {jQuery} $ele
         */
        Login: function($ele) {
            $ele.submit(function (e) {
                e.preventDefault();
                var username = $ele.find("[name=username]").val();
                var password = $ele.find("[name=password]").val();

                if (username.length == 0) {
                    alert("Must enter a username.");
                    return;
                }

                if (password.length == 0) {
                    alert("Must enter a password.");
                    return;
                }

                $.ajax({
                    type: "POST",
                    url: BaseURL + JGoExampleApp.URL.LoginSubmit,
                    data: {
                        username: username,
                        password: password
                    },
                    success: function () {
                        window.location = BaseURL + JGoExampleApp.URL.Lobby
                    },
                    /**
                     * @param {XMLHttpRequest} xhr
                     */
                    error: function (xhr) {
                        alert("Error logging in:\n" + xhr.responseText + "\nIf this problem persists, try refreshing the page.");
                    }
                });
            });
        }
    };

    var WS = {
        Type: {
            CTS: {
                SendHeartBeat: "HeartBeat",
                SendMessage: "SendMessage"
            },
            STC: {
                UserEnter: "UserEnter",
                UserExit: "UserExit",
                ReceiveMessage: "ReceiveMessage"
            }
        }
    };

    JGoExampleApp.Templates = {
        /**
         * @param {jQuery} $ele
         * @param {WebSocket} socket
         * @param {string} title
         */
        Chatroom: function ($ele, socket, title) {
            var html =
                "<div class='form-control messages'></div><br/>" +
                "<form id='chat-messages'>" +
                "<p><input type='text' class='form-control' name='message' placeholder='Type message here...' autofocus></p>" +
                "<p><input type='submit' class='btn btn-primary btn-block' value='Send'></p>" +
                "</form>";
            $ele.html(JGoExampleApp.Templates.Snippets.Panel(title, html));

            /** @type {[JGo_User]} memberList */
            var memberList = [];

            /**
             * @param {MessageEvent} msg
             */
            socket.onmessage = function (msg) {
                /** @type {JGo_WSMessage} wsMessage */
                var wsMessage = JSON.parse(msg.data);
                var html, user;

                switch (wsMessage.Type) {
                    case WS.Type.STC.ReceiveMessage:
                        /** @type {JGo_Message} wsReceiveMessage */
                        var wsReceiveMessage = JSON.parse(wsMessage.Data);
                        var $messages = $ele.find('div.messages');
                        html =
                            "<p>" +
                            "<i>[" + formatDate(wsReceiveMessage.Date) + "]</i> " +
                            "<b>" + wsReceiveMessage.User.Username + ":</b> " +
                            wsReceiveMessage.Message +
                            "</p>";
                        $messages.append(html);
                        $messages.scrollTop($messages[0].scrollHeight);
                        break;
                    case WS.Type.STC.UserEnter:
                        /** @type {JGo_User} user */
                        user = JSON.parse(wsMessage.Data);
                        memberList.push(user);
                        loadMemberList();
                        break;
                    case WS.Type.STC.UserExit:
                        /** @type {JGo_User} user */
                        user = JSON.parse(wsMessage.Data);
                        for (var i = 0; i < memberList.length; i++) {
                            if (memberList[i].Id == user.Id) {
                                memberList.splice(i, 1);
                            }
                        }
                        loadMemberList();
                }

                function loadMemberList() {
                    var $memberList = $ele.parent().find('#user-list .members');
                    if (!$memberList.length) {
                        return;
                    }
                    /**
                     * @param {JGo_User} a
                     * @param {JGo_User} b
                     */
                    function sortMemberList(a, b) {
                        return a.Username < b.Username ? -1 : 1;
                    }

                    memberList.sort(sortMemberList);
                    html = "";
                    for (var i = 0; i < memberList.length; i++) {
                        var user = memberList[i];
                        html +=
                            "<p>" +
                            "<b>" + user.Username + "</b> " +
                            "</p>";
                    }
                    var currentScroll = $memberList[0].scrollHeight;
                    $memberList.html(html);
                    $memberList.scrollTop(currentScroll);
                }
            };

            /**
             * @param {jQuery.Event} e
             */
            function chatSubmit(e) {
                e.preventDefault();
                var $msg = $ele.find('[name=message]');
                var message = $msg.val();

                if (message.length == 0) {
                    return;
                }

                /** @type {JGo_SendMessage} sendMessage */
                var sendMessage = {
                    Message: message
                };

                /** @type {JGo_WSMessage} wsMessage */
                var wsMessage = {
                    Type: WS.Type.CTS.SendMessage,
                    Data: JSON.stringify(sendMessage)
                };
                socket.send(JSON.stringify(wsMessage));
                $msg.val('');
            }

            $ele.find('#chat-messages').submit(chatSubmit);
        },
        /**
         * @param {jQuery} $ele
         * @param {WebSocket} socket
         */
        UserList: function ($ele) {
            var html = "<div id='member-list' class='form-control members'></div>";
            $ele.html(JGoExampleApp.Templates.Snippets.Panel("User List", html));
        },
        Snippets: {
            /**
             * @param {string} title
             * @param {string} html
             * @return {string}
             */
            Panel: function (title, html) {
                html =
                    "<div class='panel panel-default'>" +
                    "<div class='panel-heading'><h3 class='panel-title'>" + title + "</h3></div>" +
                    "<div class='panel-body'>" +
                    html +
                    "</div>" +
                    "</div>";
                return html;
            }
        }
    };

    JGoExampleApp.Chat = {
        Socket: null,
        /**
         * @param {jQuery} $ele
         */
        LobbyChat: function($ele) {
            var socket = getSocket(JGoExampleApp.URL.Chat);
            JGoExampleApp.Templates.Chatroom($ele, socket, "Lobby Chat");
        },
        /**
         * @param {jQuery} $ele
         */
        UserList: function($ele) {
            JGoExampleApp.Templates.UserList($ele);
        }
    };

    /**
     * @param {string} path
     * @return {WebSocket}
     */
    function getSocket(path) {
        var loc = window.location;
        var protocol = window.location.protocol.toLowerCase() == "https:" ? "wss" : "ws";
        var socket = new WebSocket(protocol + "://" + loc.hostname + ":" + loc.port + getBaseUrl() + path);

        socket.onopen = function () {
            console.log("Socket opened to: " + path);
        };

        socket.onclose = function () {
            console.log("Socket closed to: " + path);
        };

        setInterval(function () {
            var wsMessage = {
                Type: WS.Type.CTS.SendHeartBeat
            };
            socket.send(JSON.stringify(wsMessage));
        }, 15000);

        return socket;
    }

    function getBaseUrl() {
        return $("base").attr("href");
    }

    function formatDate(epochTimestamp) {
        return (new Date(epochTimestamp * 1000)).format("yyyy-MM-dd hh:mm:ss");
    }

    JGoExampleApp.URL = {
        Chat: "chat",
        Lobby: "lobby",
        LoginSubmit: "login-submit",
        SignupSubmit: "signup-submit"
    };

});

//author: meizz
Date.prototype.format = function (format) {
    var o = {
        "M+": this.getMonth() + 1, //month
        "d+": this.getDate(),      //day
        "h+": this.getHours(),     //hour
        "m+": this.getMinutes(),   //minute
        "s+": this.getSeconds(),   //second
        "q+": Math.floor((this.getMonth() + 3) / 3), //quarter
        "S": this.getMilliseconds() //millisecond
    };

    if (/(y+)/.test(format)) format = format.replace(RegExp.$1,
        (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)if (new RegExp("(" + k + ")").test(format))
        format = format.replace(RegExp.$1,
            RegExp.$1.length == 1 ? o[k] :
                ("00" + o[k]).substr(("" + o[k]).length));
    return format;
};

/**
 * @typedef {{
 *   Id: int
 *   Username: string
 * }} JGo_User
 */

/**
 * @typedef {{
 *   Type: string
 *   Data: string
 * }} JGo_WSMessage
 */

/**
 * @typedef {{
 *   Message: string
 * }} JGo_SendMessage
 */

/**
 * @typedef {{
 *   Date: int
 *   Message: string
 *   Chatroom: string
 *   User: JGo_User
 *   UserId: int
 * }} JGo_Message
 */

/**
 * https://developer.mozilla.org/en-US/docs/Web/API/MessageEvent
 * @typedef {{
 *   data: string
 * }} MessageEvent
 */
