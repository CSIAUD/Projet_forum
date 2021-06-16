var tabs = document.querySelectorAll('.info-box li');

tabs[0].click()

// var tabs = document.querySelectorAll('.info-box li a');
// var panels = document.querySelectorAll('.info-box article');

// for(i = 0; i < tabs.length; i++) {
//     var tab = tabs[i];
//     setTabHandler(tab, i);
// }

// function setTabHandler(tab, tabPos) {
//     tab.onclick = function() {
//         for(i = 0; i < tabs.length; i++) {
//             tabs[i].className = '';
//         }
//         tab.className = 'active';
//         for(i = 0; i < panels.length; i++) {
//             panels[i].className = '';
//         }
//         panels[tabPos].className = 'active-panel';
//     }
// }

function afficher(ev){
    let position = 0
    let target = ev.target
    let siblings = target.parentElement.children
    let panels = (document.getElementById("panels")).children
    let targetId = (target.id).slice(4)
    for(let panel of panels){
        if(panel.id == targetId){
            break
        }
        position++
    }
    for(let i=0;i<panels.length;i++){
        if(i==position){
            siblings[i].classList.add("active")
            panels[i].style.display = "block"
        }else{
            siblings[i].classList.remove("active")
            panels[i].style.display = "none"
        }
    }
}