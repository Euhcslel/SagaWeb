// Функция для отображения плавильного item у tab
function showItem(e, item) {
  var tabContent = document.getElementsByClassName("tabItem");
  tabContent.forEach((content) => {
    content.style.display = "none";
  });

  var tabLinks = document.getElementsByClassName("tablinks");
  tabLinks.forEach((tabLink) => {
    tabLink.classList.remove("active");
  });
  if (item === "mainConfig") {
    document.getElementById(item).style.display = "flex";
    return;
  }
  document.getElementById(item).style.display = "block";
  e.currentTarget.classList.add("active");
}
