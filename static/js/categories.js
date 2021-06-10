let popupWeek = document.getElementById("popupWeek")
let popupMonth = document.getElementById("popupMonth")
let popupAll = document.getElementById("popupAll")

let statWeek = document.querySelector("#statWeek > #WeekStats")
let statMonth = document.querySelector("#statMonth > #MonthStats")
let statAll = document.querySelector("#statAll > #AllStats")

let urlRegex = new RegExp("([0-9a-z:/\._]+)","gi")
let urlTestRegex = new RegExp("\#","g")


window.addEventListener('load', () => {
    let stats = getStats()
    console.log(stats)
    if ((window.location.href).search(urlTestRegex) > 0){
        let url = (urlRegex.exec(window.location.href))[0]
        window.location.href = url
        document.body.style.overflow = "auto"
    }
})

popupWeek.addEventListener('click', () => {
    window.location.href = (urlRegex.exec(window.location.href))[1]
    document.body.style.overflow = "auto"
})
popupMonth.addEventListener('click', () => {
    window.location.href = (urlRegex.exec(window.location.href))[1]
    document.body.style.overflow = "auto"
})
popupAll.addEventListener('click', () => {
    window.location.href = (urlRegex.exec(window.location.href))[1]
    document.body.style.overflow = "auto"
})

statWeek.addEventListener('click', () => {
    popupWeek.style.display = "block"
    setTimeout(function() {
        window.location.href = window.location.href + "#popupWeek"
    },10);
    document.body.style.overflow = "hidden"
})

statMonth.addEventListener('click', () => {
    setTimeout(function() {
        window.location.href = window.location.href + "#popupMonth"
    },10);
    document.body.style.overflow = "hidden"
})

statAll.addEventListener('click', () => {
    setTimeout(function() {
            window.location.href = window.location.href + "#popupAll"
        },10);
    document.body.style.overflow = "hidden"
})
function cut(txt){
    let temp = ""
    for (let i=txt.length; i--; i>=0)temp += txt[i]
    return txt.slice(0,(0 - temp.indexOf('/')))
}


function getStats(){
    var httpRequest
    let btn = document.getElementById("ajaxButton")
    btn.addEventListener('click', () => makeRequest());

}
function makeRequest() {
    httpRequest = new XMLHttpRequest();

    if (!httpRequest) {
        alert('Abandon :( Impossible de créer une instance de XMLHTTP');
        return false;
    }
    httpRequest.onreadystatechange = alertContents;
    httpRequest.open('GET', 'getStats');
    httpRequest.send();
}

function alertContents() {
    if (httpRequest.readyState === XMLHttpRequest.DONE) {
        if (httpRequest.status === 200) {
            console.log("ok")
            alert(httpRequest.responseText);
            return httpRequest.responseText
        } else {
            alert('Il y a eu un problème avec la requête.');
        }
    }
}