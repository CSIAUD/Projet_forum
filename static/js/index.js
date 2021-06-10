function vote(ev){
    let elem = ev.originalTarget
    let userIdCookie = getCookie("Session")
    let postId = ((((elem.parentElement).parentElement).parentElement).id).slice(4)
    let vote = elem.attributes.alt.nodeValue =="like" ? "1" : "-1"
    let url = "vote?Session="+userIdCookie+"&PostId="+postId+"&Vote="+vote
    console.log(url)
}