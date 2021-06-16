// Gestion de la pop up de login

var btnPopup = document.getElementById('login');
var overlay = document.getElementById('content');
var btnClose = document.getElementById('btnClose');

btnPopup.addEventListener('click',openlogin);
btnClose.addEventListener('click', closelogin);

function openlogin() {
    overlay.style.display='block';
}

function closelogin() {
    overlay.style.display='none';
}