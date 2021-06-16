let input = document.querySelector('#search')
var data = document.querySelectorAll('*[id^=username_]')

function searchUser() {
    let val = input.value.toLowerCase()
    for(let user of data){
        let txt = (user.innerText).toLowerCase()
        if(txt.indexOf(val) != -1){
            user.style.display = "unset"
        }else{
            user.style.display = "none"
        }
    }
}


document.getElementById("roleForm").addEventListener('submit', (ev) => {
    ev.preventDefault()
    let val = input.value

    for(let user of data){
        let txt = user.innerText
        if(val === txt){
            ev.target.submit()
        }
    }
})

function toSearch(ev){
    let target = ev.target
    document.getElementById("search").value = target.innerText
}