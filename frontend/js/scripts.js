function friendlyError(code, message) {
    switch(code) {
        case 1101:
            return "Ooops! The item already exists in the system. If you want to add another, please use another ID for it or edit the item to have a higher quantity.";
        case 1001:
            return "I'm afraid I can't let you do that Dave. Looks like the current user you are logged in with (if any) does not have permissions to perform this action.";
        default:
            return "Ooops! Theres an unexpected error \nCode: " + code + "\nMessage: " + message;
    }
}