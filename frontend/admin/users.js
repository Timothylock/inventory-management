function getUsers() {
    $.ajax({
        url: '/api/users',
        type: 'GET',
        success: function (response) {
            var result = "";
            for (var i = 0; i < response.length; i++) {
                (function (i) {
                    if (response[i].username === "System") {
                        result += "<tr><td>" + response[i].username + "</td><td>" + response[i].email + "</td><td>" + response[i].isSysAdmin + "</td><td><button type='button' class='btn btn-primary' disabled>Edit</button> <button type='button' class='btn btn-danger' disabled>Delete</button></td></tr>"
                    } else {
                        result += "<tr><td>" + response[i].username + "</td><td>" + response[i].email + "</td><td>" + response[i].isSysAdmin + "</td><td><button onclick='editUser(\"" + response[i].username + "\",\"" + response[i].real_name + "\",\"" + response[i].access_level + "\")' type='button' class='btn btn-primary' disabled>Edit</button> <button onclick='deleteUser(\"" + response[i].username + "\")' type='button' class='btn btn-danger' disabled>Delete</button></td></tr>"
                    }
                })(i);
            }
            $("#logbody").html(result);
        },
        error: function (response) {
            $("#logbody").html("<tr><td>Error</td><td>Error</td><td>Error</td><td>Error</td></tr>");
            alert("Cannot retrieve users - " + JSON.parse(response.responseText).details);
        }
    });
}

function deleteUser(username) {
    if (confirm("You're about to delete " + username + ". This cannot be undone and will be logged! Do you want to proceed?")) {
        $.ajax({
            url: '/api/admin/users?username=' + username,
            type: 'DELETE',
            success: function () {
                getUsers();
                alert("User " + username + " deleted");
            },
            error: function (response) {
                alert("Could not delete " + username + ". " + response.responseText);
            }
        });
    }
}

function editUser(username, realname, accesslevel) {
    window.location.href = "addUser.html#username=" + username + "&fullname=" + realname + "&accesslevel=" + accesslevel;
}

function submitUser() {
    var username = $('#username').val();
    var password = $('#password').val();
    var email = $('#email').val();
    var accesslevel = $('#accesslevel').val();

    console.log("about to submit");

    console.log(JSON.stringify({username: username, password: password, email: email, isSysAdmin: accesslevel}));
    $.ajax({
        url: '/api/user/add',
        type: 'POST',
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify({username: username, password: password, email: email, isSysAdmin: accesslevel}),
        success: function () {
            alert("User " + username + " added/changed");
            window.location.href = "users.html";
        },
        error: function (response) {
            alert(response.responseText);
            window.location.href = "users.html";
        }
    });
}