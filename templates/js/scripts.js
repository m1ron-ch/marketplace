function ShowMenu(id){
    document.getElementById(id).style.display = "block"
    document.getElementById("bg-menu").style.display = "block"
}

function HideMenu(id){
    document.getElementById(id).style.display = "none"
    document.getElementById("bg-menu").style.display = "none"
}

function ShowLoginMenu(){
    document.getElementById("loginMenu").style.display = "block"
    document.getElementById("signInMenu").style.display = "none"
}

function ShowSignInMenu(){
    document.getElementById("loginMenu").style.display = "none"
    document.getElementById("signInMenu").style.display = "block"
}

function EyePassword(eyeID){
    var x = document.getElementById("inputPassword");
    var boxes = document.getElementById("eyePassword");
    if (x.type === "password") {
        x.type = "text";
        boxes.className = "fa fa-eye";
    } else {
        x.type = "password";
        boxes.className = "fa fa-eye-slash";
    }
}

function EditFindCountry(){
    var selectBox = document.getElementById("countrySelectBox");
    var selectedValue = selectBox.options[selectBox.selectedIndex].value;

    var countryText = document.getElementById("whereFind");
    if (selectedValue == "AllOfBelarus"){
        countryText.innerHTML = "All ads in Belarus";
    } 
    else {
        countryText.innerHTML = "All ads in " + selectedValue;
    }
}

function Favorite(obj, adID) {
    if (obj.parentNode.className == "favorite-icon"){
        SelectFavorite(obj)

        const Http = new XMLHttpRequest();
        const url = "http://127.0.0.1:8000/add_favorites/" + adID;
        Http.open("POST", url);
        Http.send();
    
        location.reload();
    }
    else {
        DeselectFavorite(obj)

        const Http = new XMLHttpRequest();
        const url = "http://127.0.0.1:8000/delete_favorite/" + adID;
        Http.open("DELETE", url);
        Http.send();
    
        location.reload();
    }
    return false;
}

function SelectFavorite(obj) {
    var start = obj.parentNode
    start.classList.remove("favorite-icon")
    start.classList.add("favorite-icon-selected")
}

function DeselectFavorite(obj) {
    var start = obj.parentNode
    start.classList.add("favorite-icon")
    start.classList.remove("favorite-icon-selected")
}
