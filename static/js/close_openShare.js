var sharebtn = document.getElementById('sharebtn');
var share = document.getElementById('share');
var shareClose = document.getElementById('shareclose');

sharebtn.addEventListener('click',openshare);
shareClose.addEventListener('click', closeshare);

function openshare() {
    share.style.display='block';
}

function closeshare() {
    share.style.display='none';
}