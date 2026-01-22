// Переменная, хранящая начальную стоимость ворот по ширине и высоте (без выбранной комлектации)
let dealerSizePrice = 0;
let clientSizePrice = 0;
// Переменная, хранящая полную стоимость конфигурации ворот
let dealerGatePrice = 0;
let clientGatePrice = 0;

// Функция, высчитывающая полную стоимость заказа
function updateTotalPrice() {
  var totalDealerPriceElement = document.querySelector(".totalDealerPrice");
  var totalClientPriceElement = document.querySelector(".totalClientPrice");

  var liftTypeSelect = document.getElementsByName("liftType")[0];
  var selectedLiftType = liftTypeSelect.options[liftTypeSelect.selectedIndex];
  var liftTypeDealerMarkup =
    (dealerSizePrice * parseFloat(selectedLiftType.dataset.markup)) / 100;
  var liftTypeClientMarkup =
    (clientSizePrice * parseFloat(selectedLiftType.dataset.retail)) / 100;

  var cycleAmountSelect = document.getElementsByName("cycleAmount")[0];
  var selectedCycleAmount =
    cycleAmountSelect.options[cycleAmountSelect.selectedIndex];
  var cycleAmountDealerMarkup =
    (dealerSizePrice * parseFloat(selectedCycleAmount.dataset.markup)) / 100;
  var cycleAmountClientMarkup =
    (clientSizePrice * parseFloat(selectedCycleAmount.dataset.retail)) / 100;

  var drivesSelect = document.getElementsByName("drive")[0];
  var selectedDrive = drivesSelect.options[drivesSelect.selectedIndex];

  var dealerGatePrice =
    dealerSizePrice +
    liftTypeDealerMarkup +
    cycleAmountDealerMarkup +
    parseFloat(selectedDrive.dataset.price);

  var clientGatePrice =
    clientSizePrice +
    liftTypeClientMarkup +
    cycleAmountClientMarkup +
    parseFloat(selectedDrive.dataset.retail);

  totalDealerPriceElement.textContent = dealerGatePrice;
  totalClientPriceElement.textContent = clientGatePrice;
}

// Функция, определяющая начальную стоимость ворот по ширине и высоте
async function updateSizePrice() {
  const url = new URL(window.location.href);
  const gateType = url.searchParams.get("gateType");

  if (gateType === null) {
    dealerSizePrice = 0;
    clientSizePrice = 0;
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
    dealerSizePrice = data.dealer_price;
    clientSizePrice = data.client_price;
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
    colorIn: document.querySelector('input[name="colorIn"]:checked').dataset.id,
    colorOut: document.querySelector('input[name="colorOut"]:checked').dataset
      .id,
    drive:
      document.querySelector('[name="drive"]').selectedOptions[0].dataset.id,
    cycleAmount: document.querySelector('[name="cycleAmount"]')
      .selectedOptions[0].dataset.id,
    gatePrice: dealerGatePrice,
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
