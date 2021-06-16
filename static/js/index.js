function vote(ev){
    let elem = ev.target
    let postId = ((((elem.parentElement).parentElement).parentElement).id).slice(4)
    let vote = elem.attributes.alt.nodeValue =="like" ? "1" : "0"
    let url = "vote?PostId="+postId+"&Vote="+vote
    window.location.href = url
    console.log(url)
}
