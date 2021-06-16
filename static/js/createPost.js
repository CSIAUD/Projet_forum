function recupPost(ev){
    let form = ev.target
    ev.preventDefault()
    let content = document.getElementById("content")
    content.value = document.getElementById("inputPost").innerText
    form.submit()
}

function recupContent(ev){
    ev.preventDefault()
    let target = ev.target
    document.getElementById("categorie").value = target.value
}