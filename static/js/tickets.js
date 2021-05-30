
function cut(txt){
    let temp = ""
    for (let i=txt.length; i--; i>=0)temp += txt[i]
    return txt.slice(0,(0 - temp.indexOf('/')))
}