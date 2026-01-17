// Функция для отображения плавильного item у tab
function showItem(evt, item) {
  var i, tabcontent, tablinks;
  tabcontent = document.getElementsByClassName("tabItem");
  for (i = 0; i < tabcontent.length; i++) {
    tabcontent[i].style.display = "none";
  }
  /*tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
      tablinks[i].className = tablinks[i].className.replace(" active", "");
    }*/
  document.getElementById(item).style.display = "block";
  //evt.currentTarget.className += " active";
}
