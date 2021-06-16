let modos = document.getElementById("modo-container")

modos.addEventListener('click', (ev) => {
    if(ev.target.tagName == "LI" || ev.target.tagName == "FORM"){
        let id = ev.target.id
        for(let modo of modos.children){
            let modoId=modo.firstElementChild.id
            if(modoId==id){
                modo.firstElementChild.classList.add("active")
            }else{
                modo.firstElementChild.classList.remove("active")
            }
        }
    }
})