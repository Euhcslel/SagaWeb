// Объект, хранящий номер стату заказа для protobuf
export const OrderStatus = {
  pending: 0,
  confirmed: 1,
  paid: 2,
  done: 3,
  cancelled: 4,
};

// Объект, хранящий номер типа ворот для protobuf
export const gateTypeProto = {
  industrial: 0,
  residential: 1,
};

// Функция, удаляющая элемент из списка товаров или дополнительных опций
export function removeItem(element, itemContainer) {
  const item = element.closest(itemContainer);
  if (!item) return;

  const additionalOptions = item.closest(".additional-options");
  const productList = item.closest("#product-list");

  item.remove();

  if (additionalOptions) {
    updateOptionsPrice(additionalOptions);
    return;
  }

  if (productList) {
    updateProductsPrice();
  }
}

// Функция, которая добавляет новую опцию в список.
// Клонирует шаблон option-template
export function addGateOption(element) {
  const additionalOptions = element.closest(".additional-options");
  const optionList = additionalOptions.querySelector(".option-list");
  const optionTemplate = document.getElementById("option-template");

  const optionClone = optionTemplate.content.cloneNode(true);
  optionList.append(optionClone);

  updateOptionsPrice(element);
}

// Функция, которая добавляет новый товар в список.
// Клонирует шаблон product-template
export function addProduct() {
  const productList = document.getElementById("product-list");
  const productTemplate = document.getElementById("product-template");

  const productClone = productTemplate.content.cloneNode(true);
  productList.append(productClone);

  updateProductsPrice();
}

// Функция, которая обновляет цену за дополнительные опции у ворот.
// Вызывается при изменении значений или состава дополнительных опций
export function updateOptionsPrice(element) {
  let optionsRetailPrice = 0;
  let optionsWholesalePrice = 0;

  const additionalOptions = element.closest(".additional-options");
  const optionList = additionalOptions.querySelector(".option-list");
  const optionItems = optionList.getElementsByClassName("option-item");
  [...optionItems].forEach((optionItem) => {
    const optionSelect = optionItem.querySelector(".option");
    const selectedOption = optionSelect.options[optionSelect.selectedIndex];
    const amount = Number(optionItem.querySelector(".amount").value);
    optionsRetailPrice += Number(selectedOption.dataset.retailPrice) * amount;
    optionsWholesalePrice +=
      Number(selectedOption.dataset.wholesalePrice) * amount;
  });

  optionList.dataset.retailPrice = optionsRetailPrice;
  optionList.dataset.wholesalePrice = optionsWholesalePrice;

  updateGatePrice(element.closest(".gate-item"));
}

// Функция, которая обновляет цену за товары.
// Вызывается при изменении значений или состава товаров
export function updateProductsPrice() {
  let productsRetailPrice = 0;
  let productsWholesalePrice = 0;

  const productList = document.getElementById("product-list");
  const productItems = productList.getElementsByClassName("product-item");
  [...productItems].forEach((productItem) => {
    const productSelect = productItem.querySelector(".product");
    const selectedOption = productSelect.options[productSelect.selectedIndex];
    const amount = Number(productItem.querySelector(".amount").value);
    productsRetailPrice += Number(selectedOption.dataset.retailPrice) * amount;
    productsWholesalePrice +=
      Number(selectedOption.dataset.wholesalePrice) * amount;
  });

  productList.dataset.retailPrice = productsRetailPrice;
  productList.dataset.wholesalePrice = productsWholesalePrice;

  updateOrderPrice();
}

// Функция, которая обновляет цену на текущие ворота.
// Вызывается при изменении занчений select и input, которые влияют на стоимость ворот
export function updateGatePrice(gate) {
  let retailGatePrice = 0;
  let wholesaleGatePrice = 0;

  const size = gate.querySelector(".sizes");
  const sizeRetailPrice = Number(size.dataset.retailPrice || 0);
  const sizeWholesalePrice = Number(size.dataset.wholesalePrice || 0);

  const liftType = gate.querySelector(".lift-type");
  const selectedLiftType = liftType.options[liftType.selectedIndex];
  const liftTypeRetailMarkup = Number(selectedLiftType.dataset.retailMarkup);
  const liftTypeWholesaleMarkup = Number(
    selectedLiftType.dataset.wholesaleMarkup,
  );

  const cycleAmount = gate.querySelector(".cycle-amount");
  const selectedCycleAmount = cycleAmount.options[cycleAmount.selectedIndex];
  const cycleAmountRetailMarkup = Number(
    selectedCycleAmount.dataset.retailMarkup,
  );
  const cycleAmountWholesaleMarkup = Number(
    selectedCycleAmount.dataset.wholesaleMarkup,
  );

  let driveRetailPrice = 0;
  let driveWholesalePrice = 0;
  const driveType = gate.querySelector(".drive-type");
  switch (driveType.value) {
    case "manual": {
      const chainLengthInput = gate.querySelector(".chain-length");
      const chainLength = Number(chainLengthInput.value);
      driveRetailPrice +=
        Number(chainLengthInput.dataset.chainDriveRetailPrice) +
        Number(chainLength) * Number(chainLengthInput.dataset.chainRetailPrice);
      driveWholesalePrice +=
        Number(chainLengthInput.dataset.chainDriveWholesalePrice) +
        Number(chainLength) *
          Number(chainLengthInput.dataset.chainWholesalePrice);
      break;
    }
    case "residential": {
      const drive = gate.querySelector(".drive");
      const selectedDrive = drive.options[drive.selectedIndex];
      const rail = gate.querySelector(".rail");
      const selectedRail = rail.options[rail.selectedIndex];
      driveRetailPrice +=
        Number(selectedDrive.dataset.retailPrice) +
        Number(selectedRail.dataset.retailPrice);
      driveWholesalePrice +=
        Number(selectedDrive.dataset.wholesalePrice) +
        Number(selectedRail.dataset.wholesalePrice);
      break;
    }
    case "industrial": {
      const drive = gate.querySelector(".drive");
      const selectedDrive = drive.options[drive.selectedIndex];
      driveRetailPrice += Number(selectedDrive.dataset.retailPrice);
      driveWholesalePrice += Number(selectedDrive.dataset.wholesalePrice);
      break;
    }
  }

  var optionList = gate.querySelector(".option-list");
  var optionListRetailPrice = Number(optionList.dataset.retailPrice || 0);
  var optionListWholesalePrice = Number(optionList.dataset.wholesalePrice || 0);

  retailGatePrice =
    sizeRetailPrice +
    (sizeRetailPrice * liftTypeRetailMarkup) / 100 +
    (sizeRetailPrice * cycleAmountRetailMarkup) / 100 +
    driveRetailPrice +
    optionListRetailPrice;

  wholesaleGatePrice =
    sizeWholesalePrice +
    (sizeWholesalePrice * liftTypeWholesaleMarkup) / 100 +
    (sizeWholesalePrice * cycleAmountWholesaleMarkup) / 100 +
    driveWholesalePrice +
    optionListWholesalePrice;

  gate.dataset.retailPrice = retailGatePrice;
  gate.dataset.wholesalePrice = wholesaleGatePrice;
  document.getElementById("gate-retail-price").textContent = retailGatePrice.toFixed(2);
  document.getElementById("gate-wholesale-price").textContent =
    wholesaleGatePrice.toFixed(2);
};
