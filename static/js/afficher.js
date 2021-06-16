// cette fonction permet d'afficher les mots de passes en appuyant sur l'élement image (oeil)

function Afficher(ev,id){
    let target = ev.target
    var input = document.getElementById(id); 
    console.log(input.value)
    if (input.type === "password"){ 
        input.setAttribute("type", "text")
        target.src = "./static/img/oeil non barre.png"
    }else{
        input.setAttribute("type", "password")
        target.src = "./static/img/oeil barre.png"
    } 
} 

function IndexRedirect() {
    window.location.href="index";
}