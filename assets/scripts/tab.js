// Функция для отображения плавильного item у tab
function showItem(e, item) {
  var i, tabcontent, tablinks;
  tabcontent = document.getElementsByClassName("tabItem");
  for (i = 0; i < tabcontent.length; i++) {
    tabcontent[i].style.display = "none";
  }

  tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
      tablinks[i].classList.remove("active")
    }
  if (item === "mainConfig") {
    document.getElementById(item).style.display = "flex";
    return;
  }
  document.getElementById(item).style.display = "block";
  e.currentTarget.classList.add("active");
}
