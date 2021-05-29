
var btnPopup = document.getElementById('btnPopup');
var overlay = document.getElementById('overlay');
var btnClose = document.getElementById('btnClose');
btnPopup.addEventListener('click',openModal);
btnClose.addEventListener('click', closeModal);
function openModal() {
overlay.style.display='block';
}

function closeModal() {
    overlay.style.display='none';
}