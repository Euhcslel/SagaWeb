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

  var optionsPrice = 0;
  var optionDivs = document.querySelectorAll(".option-div");
  optionDivs.forEach((div) => {
    var select = div.querySelector('select[name="option"]');
    var amountInput = div.querySelector('input[type="number"]');

    if (select && amountInput) {
      var selectedOption = select.options[select.selectedIndex];
      var price = parseFloat(selectedOption.dataset.price) || 0;
      var amount = parseInt(amountInput.value) || 0;

      optionsPrice += price * amount;
    }
  });

  var productsPrice = 0;
  var productDivs = document.querySelectorAll(".product-div");
  productDivs.forEach((div) => {
    var select = div.querySelector('select[name="product"]');
    var amountInput = div.querySelector('input[type="number"]');

    if (select && amountInput) {
      var selectedProduct = select.options[select.selectedIndex];
      var price = parseFloat(selectedProduct.dataset.price) || 0;
      var amount = parseInt(amountInput.value) || 0;

      productsPrice += price * amount;
    }
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
    gateTypeId: parseInt(document.querySelector(".gateType").value),
    width: parseInt(document.getElementById("width").value),
    height: parseInt(document.getElementById("height").value),
    liftTypeId: parseInt(
      document.querySelector('[name="liftType"]').selectedOptions[0].dataset.id,
    ),
    colorOutId: parseInt(
      document.querySelector('input[name="colorOut"]:checked').dataset.id,
    ),
    driveId: parseInt(
      document.querySelector('[name="drive"]').selectedOptions[0].dataset.id,
    ),
    cycleAmountId: parseInt(
      document.querySelector('[name="cycleAmount"]').selectedOptions[0].dataset
        .id,
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
  const root = await protobuf.load("/assets/scripts/order.proto");
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

  await fetch("/orders", {
    method: "POST",
    headers: { "Content-Type": "application/x-protobuf" },
    body: buffer,
  });
}

function deleteItem(button, divName) {
  const optionDiv = button.closest(divName);
  if (optionDiv) {
    optionDiv.remove();
    updateTotalPrice();
  }
}
