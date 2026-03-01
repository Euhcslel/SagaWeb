// Переменная, хранящая начальную стоимость ворот по ширине и высоте (без выбранной комлектации)
let sizePrice = 0;
// Переменная, хранящая полную стоимость конфигурации ворот
let gatePrice = 0;

document.addEventListener("input", function (event) {
  if (event.target.matches("input[type='number']")) {
    updateSizePrice();
  }
});

// Функция, высчитывающая полную стоимость заказа
function updateTotalPrice() {
  var totalPriceElement = document.querySelector(".totalPrice");

  var liftInput = document.querySelector('input[name="liftType"]');
  var liftMarkup = Number(liftInput.dataset.markup || 0);
  var liftMarkupInMoney = (sizePrice * liftMarkup) / 100;

  var cycleInput = document.querySelector('input[name="cycleAmount"]');
  var cycleMarkup = Number(cycleInput.dataset.markup || 0);
  var cycleMarkupInMoney = (sizePrice * cycleMarkup) / 100;

  var driveInput = document.querySelector('input[name="drive"]');
  var drivePrice = Number(driveInput.dataset.price || 0);

  var optionsPrice = 0;
  var optionDivs = document.querySelectorAll(".option-div");
  optionDivs.forEach((div) => {
    var optionInput = document.querySelector('input[name="option"]');
    var amountInput = div.querySelector('input[type="number"]');

    if (optionInput && amountInput) {
      var optionPrice = Number(optionInput.dataset.price || 0);
      var amount = parseInt(amountInput.value) || 0;

      optionsPrice += optionPrice * amount;
    }
  });

  var productsPrice = 0;
  var productDivs = document.querySelectorAll(".product-div");
  productDivs.forEach((div) => {
    var productInput = document.querySelector('input[name="product"]');
    var amountInput = div.querySelector('input[type="number"]');

    if (productInput && amountInput) {
      var productPrice = Number(productInput.dataset.price || 0);
      var amount = parseInt(amountInput.value) || 0;

      productsPrice += productPrice * amount;
    }
  });

  gatePrice =
    productsPrice +
    optionsPrice +
    sizePrice +
    liftMarkupInMoney +
    cycleMarkupInMoney +
    drivePrice;

  totalPriceElement.textContent = gatePrice.toFixed(2) + " руб.";
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

// Текущий URL адрес
const url = new URL(window.location.href);

// Тип ворот, взятый из URL адреса
const gateType = url.searchParams.get("gateType");

// Элементы формы, хранящий тип ворот
gateTypeInput = document.getElementsByClassName("gateType")[0];

// Присвоение значения типа ворот в соответствующий элемент формы
gateTypeInput.value = gateType;

function getObjectWithAmount(divName, selectName) {
  var itemsList = {};

  var divs = document.getElementsByClassName(divName);
  Array.from(divs).forEach((div) => {
    var select = div.querySelector("select[name=" + selectName + "]");
    var selectedOption = select.options[select.selectedIndex];
    var optionId = selectedOption.dataset.id;
    var amountInput = div.querySelector('input[type="number"]');
    var amount = parseInt(amountInput.value) || 0;

    itemsList[optionId] = (itemsList[optionId] || 0) + amount;
  });

  return itemsList;
}

// Функция, добавляющая конфигурацию ворот в заказ
function addInOrder() {
  var optionsList = getObjectWithAmount("option-div", "option");

  const config = {
    headroom: 100,
    gateTypeId: parseInt(document.querySelector(".gateType").value),
    width: parseInt(document.getElementById("width").value),
    height: parseInt(document.getElementById("height").value),
    liftTypeId: parseInt(
      document.querySelector('input[name="liftType"]').value,
    ),
    colorOutId: parseInt(
      document.querySelector('input[name="colorOut"]').value,
    ),
    driveId: parseInt(
      document.querySelector('input[name="drive"]').value,
    ),
    cycleAmountId: parseInt(
      document.querySelector('input[name="cycleAmount"]').value,
    ),
    options: optionsList,
    gatePrice: parseInt(gatePrice),
  };

  orderGates.push(config);
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
  const root = await protobuf.load("/assets/proto_files/order.proto");
  const OrderRequest = root.lookupType("proto.OrderRequest");

  var productsList = getObjectWithAmount("product-div", "product");
  if (orderGates.length == 0) {
    addInOrder();
  }
  console.log(orderGates, productsList);
  const payload = {
    orderGates: orderGates,
    products: productsList,
  };

  const err = OrderRequest.verify(payload);
  if (err) throw new Error(err);

  const message = OrderRequest.create(payload);
  const buffer = OrderRequest.encode(message).finish();

  const res = await fetch("/orders", {
    method: "POST",
    headers: { "Content-Type": "application/x-protobuf" },
    body: buffer,
  });

  if (res.ok) {
    window.location.href = "/orders";
  }
}

function deleteItem(button, divName) {
  const optionDiv = button.closest(divName);
  if (optionDiv) {
    optionDiv.remove();
    updateTotalPrice();
  }
}
