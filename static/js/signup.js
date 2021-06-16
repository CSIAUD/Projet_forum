let valids = document.querySelectorAll("#verifs > label > span")
window.addEventListener('load', () =>{
    let elem = document.getElementById("verifs")
    if (elem.classList.contains("display")) toggle("verifs")
})

function specs(ev){
    let count = 0
    let upperRegex = new RegExp('([A-Z])')
    let lowerRegex = new RegExp('([a-z])')
    let nbRegex = new RegExp('([0-9])')
    let specialRegex = new RegExp('([!@#$%^&*])')
    let value = ev.target.value
    let length = value.length
    console.log(length)
    console.log(value)
    if(length>=8){
        valids[0].classList.add('ok')
    }else{
        valids[0].classList.remove('ok')
        count++
    }
    
    if(value.match(upperRegex)){
        valids[1].classList.add('ok')
    }else{
        valids[1].classList.remove('ok')
        count++
    }
    
    if(value.match(lowerRegex)){
        valids[2].classList.add('ok')
    }else{
        valids[2].classList.remove('ok')
        count++
    }
    
    if(value.match(nbRegex)){
        valids[3].classList.add('ok')
    }else{
        valids[3].classList.remove('ok')
        count++
    }
    
    if(value.match(specialRegex)){
        valids[4].classList.add('ok')
    }else{
        valids[4].classList.remove('ok')
        count++
    }
    console.log(value)
    return count==0
}

function toggle(id){
    document.getElementById(id).classList.toggle("display")
}