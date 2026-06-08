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
    const itemRetail = Number(selectedOption.dataset.retailPrice) * amount;
    const itemWholesale = Number(selectedOption.dataset.wholesalePrice) * amount;
    optionsRetailPrice += itemRetail;
    optionsWholesalePrice += itemWholesale;

    const priceDisplay = optionItem.querySelector(".option-price-display");
    if (priceDisplay) {
      priceDisplay.textContent = `Клиент: ${fmtPrice(itemRetail)} / Дилер: ${fmtPrice(itemWholesale)}`;
    }
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
  let driveOnlyRetailPrice = 0;
  let driveOnlyWholesalePrice = 0;
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
      driveOnlyRetailPrice = driveRetailPrice;
      driveOnlyWholesalePrice = driveWholesalePrice;
      break;
    }
    case "residential": {
      const drive = gate.querySelector(".drive");
      const selectedDrive = drive?.options[drive.selectedIndex];
      const rail = gate.querySelector(".rail");
      const selectedRail = rail?.options[rail.selectedIndex];
      if (selectedDrive) {
        driveOnlyRetailPrice = Number(selectedDrive.dataset.retailPrice);
        driveOnlyWholesalePrice = Number(selectedDrive.dataset.wholesalePrice);
      }
      if (selectedRail) {
        driveRetailPrice += driveOnlyRetailPrice + Number(selectedRail.dataset.retailPrice);
        driveWholesalePrice += driveOnlyWholesalePrice + Number(selectedRail.dataset.wholesalePrice);
      } else {
        driveRetailPrice += driveOnlyRetailPrice;
        driveWholesalePrice += driveOnlyWholesalePrice;
      }
      break;
    }
    case "industrial": {
      const drive = gate.querySelector(".drive");
      const selectedDrive = drive?.options[drive.selectedIndex];
      if (selectedDrive) {
        driveOnlyRetailPrice = Number(selectedDrive.dataset.retailPrice);
        driveOnlyWholesalePrice = Number(selectedDrive.dataset.wholesalePrice);
        driveRetailPrice += driveOnlyRetailPrice;
        driveWholesalePrice += driveOnlyWholesalePrice;
      }
      break;
    }
  }

  let optionList = gate.querySelector(".option-list");
  let optionListRetailPrice = Number(optionList.dataset.retailPrice || 0);
  let optionListWholesalePrice = Number(optionList.dataset.wholesalePrice || 0);

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
  document.getElementById("gate-retail-price").textContent =
    fmtPrice(retailGatePrice);
  document.getElementById("gate-wholesale-price").textContent =
    fmtPrice(wholesaleGatePrice);

  const liftEl = gate.querySelector(".lift-markup-display");
  if (liftEl) liftEl.textContent =
    `Наценка: клиент +${liftTypeRetailMarkup}% / дилер +${liftTypeWholesaleMarkup}%`;

  const cycleEl = gate.querySelector(".cycle-markup-display");
  if (cycleEl) cycleEl.textContent =
    `Наценка: клиент +${cycleAmountRetailMarkup}% / дилер +${cycleAmountWholesaleMarkup}%`;

  const driveEl = gate.querySelector(".drive-price-display");
  if (driveEl) driveEl.textContent =
    `Клиент: ${fmtPrice(driveOnlyRetailPrice)} / Дилер: ${fmtPrice(driveOnlyWholesalePrice)}`;

  const railEl = gate.querySelector(".rail-price-display");
  if (railEl && driveType.value === "residential") {
    const rail = gate.querySelector(".rail");
    const selectedRail = rail.options[rail.selectedIndex];
    const railRetail = Number(selectedRail.dataset.retailPrice);
    const railWholesale = Number(selectedRail.dataset.wholesalePrice);
    railEl.textContent = `Клиент: ${fmtPrice(railRetail)} / Дилер: ${fmtPrice(railWholesale)}`;
  }

  const chainEl = gate.querySelector(".chain-price-display");
  if (chainEl && driveType.value === "manual") {
    chainEl.textContent = `Клиент: ${fmtPrice(driveRetailPrice)} / Дилер: ${fmtPrice(driveWholesalePrice)}`;
  }
}

export function fmtPrice(price) {
  const roundedPrice = Math.round(price * 100) / 100;
  const [intPart, fracPart] = roundedPrice.toFixed(2).split(".");
  const formattedInt = intPart.replace(/\B(?=(\d{3})+(?!\d))/g, " ");

  return `${formattedInt},${fracPart} ₽`;
}
