// Переменная, хранящая начальную стоимость ворот по ширине и высоте (без выбранной комлектации)
let sizePrice = 0;
// Переменная, хранящая полную стоимость конфигурации ворот
let gatePrice = 0;

// Функция, высчитывающая полную стоимость заказа
function updateTotalPrice() {
  var totalPriceElement = document.querySelector(".totalPrice");
  var liftTypeSelect = document.getElementsByName("liftType")[0];
  var selectedLiftType = liftTypeSelect.options[liftTypeSelect.selectedIndex];
  var cycleAmountSelect = document.getElementsByName("cycleAmount")[0];
  var selectedCycleAmount =
    cycleAmountSelect.options[cycleAmountSelect.selectedIndex];
  var drivesSelect = document.getElementsByName("drive")[0];
  var selectedDrive = drivesSelect.options[drivesSelect.selectedIndex];

  var liftTypeMarkup =
    (sizePrice * parseFloat(selectedLiftType.dataset.markup)) / 100;
  var cycleAmountMarkup =
    (sizePrice * parseFloat(selectedCycleAmount.dataset.markup)) / 100;

  gatePrice =
    sizePrice +
    liftTypeMarkup +
    cycleAmountMarkup +
    parseFloat(selectedDrive.dataset.price);

  totalPriceElement.textContent = gatePrice;
}

// Функция, определяющая начальную стоимость ворот по ширине и высоте
async function updateSizePrice() {
  const url = new URL(window.location.href);
  const gateType = url.searchParams.get("gateType");

  if (gateType === null) {
    sizePrice = 0;
    updateTotalPrice();
    return;
  }

  var widthInput = document.getElementById("width");
  var heightInput = document.getElementById("height");

  var width = parseInt(widthInput.value) || 0;
  var height = parseInt(heightInput.value) || 0;

  if (width >= widthInput.min && height >= heightInput.min) {
    const response = await fetch(
      `/sizes?width=${width}&height=${height}&gateType=${gateType}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      },
    );

    const data = await response.json();
    sizePrice = data.price;
    updateTotalPrice();
  }
}

// Обновление начальной стоимости ворот после загрузки контента
document.addEventListener("DOMContentLoaded", function () {
  updateSizePrice();
});

// Массив для хранения объектов с конфигурацией ворот
var orderGates = [];

// Объект, хранящий выбранные товары (автоматику)
var products = {};

// Текущий URL адрес
const url = new URL(window.location.href);

// Тип ворот, взятый из URL адреса
const gateType = url.searchParams.get("gateType");

// Элементы формы, хранящий тип ворот
gateTypeInput = document.getElementsByClassName("gateType")[0];

// Присвоение значения типа ворот в соответствующий элемент формы
gateTypeInput.value = gateType;

// Функция, добавляющая конфигурацию ворот в заказ
function addInOrder() {
  const config = {
    gateType: document.querySelector(".gateType").value,
    width: document.getElementById("width").value,
    height: document.getElementById("height").value,
    liftType:
      document.querySelector('[name="liftType"]').selectedOptions[0].dataset.id,
    colorIn:
      document.querySelector('[name="colorIn"]').selectedOptions[0].dataset.id,
    colorOut:
      document.querySelector('[name="colorOut"]').selectedOptions[0].dataset.id,
    drive:
      document.querySelector('[name="drive"]').selectedOptions[0].dataset.id,
    cycleAmount: document.querySelector('[name="cycleAmount"]')
      .selectedOptions[0].dataset.id,
    gatePrice: gatePrice,
  };

  orderGates.push(config);
}

// Функция, добавляющая дополнительную комплектацию к конфигурации ворот
function addGateOption() {}

// Функция, добавляющая товар к текущему заказу
function addProductInOrder() {}

// Функция для оформления заказа
async function placeOrder() {
  const fd = new FormData();
  //const form = document.getElementById("myForm");

  fd.append("orderGates", JSON.stringify(orderGates));
  fd.append("products", JSON.stringify(products));

  var response = await fetch("/orders", {
    method: "POST",
    body: fd,
  });
  if (response.ok && response.redirected) {
    window.location.href = "/orders";
  }
}
