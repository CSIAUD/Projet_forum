function Afficher(id)
{ 
    var input = document.getElementById(id); 
    if (input.type === "password")
        { 
            input.setAttribute("type", "text")
        } 
    else
        { 
            input.setAttribute("type", "password")
        } 
} 