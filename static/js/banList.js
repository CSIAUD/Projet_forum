// On rend les éléments visibles (liste des ban) ou non 
let bans = document.getElementsByClassName("userBan")
for (ban of bans){
    let pos = ban.children[0].getBoundingClientRect()
    ban.children[1].style.left = pos.right + "px"
} 

for (user of bans){
    user.addEventListener('click', (ev) => {
        check()
        let list = ev.target.parentElement.children[1]
        
        if(list.classList.contains('display')){
            list.classList.remove("visible")
            
            setTimeout(function() {
                list.classList.remove("display")
            },550);
        }else{
            list.classList.add("display")
            setTimeout(function() {
                list.classList.add("visible")
            },10);
        }
    })
}

function check(){
    for (user of bans){
        let list = user.children[1]

        if(list.classList.contains('display')){
            list.classList.remove("visible")
            
            setTimeout(function() {
                list.classList.remove("display")
            },550);
        }
    }
}

function cut(txt){
    let temp = ""
    for (let i=txt.length; i--; i>=0)temp += txt[i]
    return txt.slice(0,(0 - temp.indexOf('/')))
}