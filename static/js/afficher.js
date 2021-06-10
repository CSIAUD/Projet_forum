function Afficher(ev,id){
    let target = ev.target
    var input = document.getElementById(id); 
    if (input.type === "password"){ 
        input.setAttribute("type", "text")
        target.src = "./static/img/oeil non barre.png"
    }else{
        input.setAttribute("type", "password")
        target.src = "./static/img/oeil barre.png"
    } 
} 