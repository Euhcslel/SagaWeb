// Переменная, хранящая начальную стоимость ворот по ширине и высоте (без выбранной комлектации)
let sizePrice = 0;
// Переменная, хранящая полную стоимость конфигурации ворот
let gatePrice = 0;

var optionsPrice = 0;
var productsPrice = 0;

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

  optionsPrice = 0;
  var optionSelects = document.getElementsByName("option");
  optionSelects.forEach((select) => {
    var selectedOption = select.options[select.selectedIndex];
    var price = selectedOption.dataset.price;
    optionsPrice += parseInt(price);
  });

  productsPrice = 0;
  var productSelects = document.getElementsByName("product");
  productSelects.forEach((select) => {
    var selectedProduct = select.options[select.selectedIndex];
    var price = selectedProduct.dataset.price;
    productsPrice += parseInt(price);
  });

  gatePrice =
    productsPrice +
    optionsPrice +
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

// Массив, хранящий выбранные товары (автоматику)
var productsList = [];

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
  var optionsList = [];
  var optionSelects = document.getElementsByName("option");
  optionSelects.forEach((select) => {
    var selectedOption = select.options[select.selectedIndex];
    optionsList.push(selectedOption.dataset.id);
  });

  const config = {
    gateType: document.querySelector(".gateType").value,
    width: document.getElementById("width").value,
    height: document.getElementById("height").value,
    liftType:
      document.querySelector('[name="liftType"]').selectedOptions[0].dataset.id,
    colorIn: document.querySelector('input[name="colorIn"]:checked').dataset.id,
    colorOut: document.querySelector('input[name="colorOut"]:checked').dataset
      .id,
    drive:
      document.querySelector('[name="drive"]').selectedOptions[0].dataset.id,
    cycleAmount: document.querySelector('[name="cycleAmount"]')
      .selectedOptions[0].dataset.id,
    options: optionsList,
    gatePrice: gatePrice,
  };

  orderGates.push(config);
  console.log(orderGates)
}

// Функция, добавляющая дополнительную комплектацию к конфигурации ворот
function addGateOption() {
  const template = document.getElementById("optionTemplate");
  const clone = template.content.cloneNode(true);
  document.getElementById("optionList").appendChild(clone);
  updateTotalPrice();
}

// Функция, добавляющая товар к текущему заказу
function addProductInOrder() {
  const template = document.getElementById("productTemplate");
  const clone = template.content.cloneNode(true);
  document.getElementById("productList").appendChild(clone);
  updateTotalPrice();
}

// Функция для оформления заказа
async function placeOrder() {
  var productSelects = document.getElementsByName("product");
  productSelects.forEach((select) => {
    var selectedProduct = select.options[select.selectedIndex];
    productsList.push(selectedProduct.dataset.id);
  });

  if (orderGates.length == 0) {
    addInOrder();
  }
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
