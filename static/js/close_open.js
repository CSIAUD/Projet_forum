var btnPopup = document.getElementById('login');
var overlay = document.getElementById('content');
var btnClose = document.getElementById('btnClose');
btnPopup.addEventListener('click',openModal);
btnClose.addEventListener('click', closeModal);
function openModal() {
overlay.style.display='block';
}

function closeModal() {
    overlay.style.display='none';
}