// Proto-схемы
let Proto = {
  SizePrice: null,
  OrderRequest: null,
};

// Функция, инициализирующая proto-схемы
async function initProtobuf() {
  const pricesRoot = await protobuf.load("/api/proto/prices.proto");
  Proto.SizePrice = pricesRoot.lookupType("proto.SizePrice");

  const orderRoot = await protobuf.load("/api/proto/order.proto");
  Proto.OrderRequest = orderRoot.lookupType("proto.OrderRequest");
}

// Функция, удаляющая вкладку ворот и соответствующий шаблон формы
window.removeGateElement = (event) => {
  const gateTabs = document.getElementById("gate-tabs");
  const gateList = document.getElementById("gate-list");

  const tabElement = event.currentTarget.closest("#gate-tabs > li");
  const tabIndex = Array.from(gateTabs.children).indexOf(tabElement);

  tabElement.remove();
  gateList.children[tabIndex]?.remove();

  const firstInput = gateTabs.querySelector(
    'li:first-child input[type="radio"]',
  );
  if (firstInput) {
    firstInput.checked = true;
  }
};

// Функция, вызывающаяся при смене текущих ворот
// необходима, чтобы корректно менять цену ворот в заказе
window.onChangeGateButton = (e) => {
  if (!e.target.matches('input[name="gates"]')) return;

  const radios = document.querySelectorAll('#gate-tabs input[name="gates"]');
  const index = [...radios].findIndex((r) => r.checked);

  const gate = document.getElementsByClassName("gate-item")[index];

  document.getElementById("gate-retail-price").textContent =
    gate?.dataset.retailPrice || 0;

  document.getElementById("gate-wholesale-price").textContent =
    gate?.dataset.wholesalePrice || 0;
};

// Функция, добавляющая новую вкладку и шаблон формы
window.addGateElement = () => {
  const gateTemplate = document.getElementById("gate-template");
  const gateList = document.getElementById("gate-list");

  const gateClone = gateTemplate.content.cloneNode(true);
  gateList.append(gateClone);

  const gateTabElementTemplate = document.getElementById("gate-tab-template");
  const gateTabs = document.getElementById("gate-tabs");

  const gateTabElementClone = gateTabElementTemplate.content.cloneNode(true);
  gateTabs.append(gateTabElementClone);

  const lastGate = gateList.lastElementChild;
  updateGateType(lastGate.querySelector(".gate-type"));
  updateGateSizePrice(lastGate.querySelector(".width"));
  updateDriveType(lastGate.querySelector(".drive-type"));
};

// Функция, которая срабатывает при смене типа ворот
// обновляет динамические значения формы, которые зависят от типа ворот
window.updateGateType = (select) => {
  const gate = select.closest(".gate-item");
  const cfg = configuration[select.value + "Configuration"];

  const width = gate.querySelector(".width");
  width.min = cfg.widthParams.MinValue;
  width.max = cfg.widthParams.MaxValue;
  width.value = cfg.widthParams.MinValue;

  const height = gate.querySelector(".height");
  height.min = cfg.heightParams.MinValue;
  height.max = cfg.heightParams.MaxValue;
  height.value = cfg.heightParams.MinValue;

  const liftType = gate.querySelector(".lift-type");
  fillSelect(liftType, cfg.liftTypes, "markup");

  const cycleAmount = gate.querySelector(".cycle-amount");
  fillSelect(cycleAmount, cfg.cycleAmounts, "markup");

  const driveType = gate.querySelector(".drive-type");
  fillSelect(driveType, cfg.driveTypes, "object");

  const drive = gate.querySelector(".drive");
  fillSelect(drive, cfg.drives, "price");
  drive.selectedIndex = 0;

  if (select.value === "residential") {
    const rail = gate.querySelector(".rail");
    fillSelect(rail, cfg.rails, "price");
    rail.selectedIndex = 0;
  }

  updateDriveType(gate.querySelector(".drive-type"));
};

// Функция, которая используется для заполнения элементов select
// в зависимости от переданного типа, создает option с нужными параметрами
function fillSelect(select, items, type) {
  select.innerHTML = "";
  if (type === "markup") {
    items.forEach((item) => {
      const option = document.createElement("option");
      option.dataset.retailMarkup = item.RetailMarkup;
      option.dataset.wholesaleMarkup = item.WholesaleMarkup;
      option.value = item.ID;
      option.textContent = item.Name;

      select.append(option);
    });
  } else if (type === "object") {
    Object.entries(items).forEach(([key, value]) => {
      const option = document.createElement("option");
      option.value = key;
      option.textContent = value;

      select.append(option);
    });
  } else if (type === "price") {
    items.forEach((item) => {
      const option = document.createElement("option");
      option.dataset.retailPrice = item.RetailPrice;
      option.dataset.wholesalePrice = item.WholesalePrice;
      option.value = item.ID;
      option.textContent = item.Name;

      select.append(option);
    });
  }
}

// Функция, срабатывающая при изменении привода
// показывает нужные поля ввода и скрывает ненужные
window.updateDriveType = (select) => {
  const gate = select.closest(".gate-item");

  const isManual = select.value === "manual";

  const chainLabel = gate.querySelector(".chain-length").closest("label");
  const driveAutoBlock = gate.querySelector(".drive-auto");
  const railLabel = gate.querySelector(".rail").closest("label");
  const gateType = gate.querySelector(".gate-type").value;

  chainLabel.hidden = !isManual;
  driveAutoBlock.hidden = isManual;
  railLabel.hidden = isManual || gateType !== "residential";

  updateGatePrice(gate);
};

// Функция, срабатывающая при изменении значений ширины и высоты
// запрашивает цены клиента и дилера у сервера и присваивает их блоку sizes
window.updateGateSizePrice = async (input) => {
  const gate = input.closest(".gate-item");
  const width = gate.querySelector(".width");
  const height = gate.querySelector(".height");
  const gateType = gate.querySelector(".gate-type");

  const sizes = input.closest(".sizes");

  const price = await fetchSizePrice(
    width.value,
    height.value,
    gateType.value.slice(0, 3),
  );

  sizes.dataset.retailPrice = price.retail;
  sizes.dataset.wholesalePrice = price.wholesale;

  updateGatePrice(gate);
};

// Функция, запрашивающая цену за размер ворот
// принимает параметры: w - ширина, h - высота, t - тип ворот
async function fetchSizePrice(w, h, t) {
  const res = await fetch(`/sizes?width=${w}&height=${h}&gateType=${t}`);

  const buf = await res.arrayBuffer();

  const msg = Proto.SizePrice.decode(new Uint8Array(buf));

  const d = Proto.SizePrice.toObject(msg, { longs: Number });

  if (d.dealer) {
    return {
      retail: d.dealer.clientPrice / 100,
      wholesale: d.dealer.dealerPrice / 100,
    };
  }

  return { retail: d.client.clientPrice / 100, wholesale: 0 };
}

// Функция, которая обновляет цену на текущие ворота.
// Вызывается при изменении занчений select и input, которые влияют на стоимость ворот
window.updateGatePrice = (gate) => {
  var retailGatePrice = 0;
  var wholesaleGatePrice = 0;

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

  var driveRetailPrice = 0;
  var driveWholesalePrice = 0;
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
  document.getElementById("gate-retail-price").textContent = retailGatePrice;
  document.getElementById("gate-wholesale-price").textContent =
    wholesaleGatePrice;

  updateOrderPrice();
};

// Функция, которая пересчитывает стоимость всего заказа
window.updateOrderPrice = () => {
  var totalRetailPrice = 0;
  var totalWholesalePrice = 0;

  const gates = document.getElementsByClassName("gate-item");
  [...gates].forEach((gate) => {
    const amount = Number(gate.querySelector(".gate-amount").value);

    const gateRetailPrice = Number(gate.dataset.retailPrice || 0);
    totalRetailPrice += gateRetailPrice * amount;

    const gateWholesalePrice = Number(gate.dataset.wholesalePrice || 0);
    totalWholesalePrice += gateWholesalePrice * amount;
  });

  const products = document.getElementById("product-list");
  totalRetailPrice += Number(products.dataset.retailPrice || 0);
  totalWholesalePrice += Number(products.dataset.wholesalePrice || 0);

  const orderRetailPriceElement = document.getElementById("order-retail-price");
  orderRetailPriceElement.textContent = totalRetailPrice;

  const orderWholesalePriceElement = document.getElementById(
    "order-wholesale-price",
  );
  orderWholesalePriceElement.textContent = totalWholesalePrice;
};

// Функция, которая обновляет цену за дополнительные опции у ворот.
// Вызывается при изменении значений или состава дополнительных опций
window.updateOptionsPrice = (element) => {
  var optionsRetailPrice = 0;
  var optionsWholesalePrice = 0;

  const additionalOptionsBlock = element.closest(".additional-options");
  const optionList = additionalOptionsBlock.querySelector(".option-list");
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
};

// Функция, которая добавляет новую опцию в список.
// Клонирует шаблон option-template
window.addGateOption = (element) => {
  const additionalOptionsBlock = element.closest(".additional-options");
  const optionList = additionalOptionsBlock.querySelector(".option-list");
  const optionTemplate = document.getElementById("option-template");

  const optionClone = optionTemplate.content.cloneNode(true);
  optionList.append(optionClone);

  updateOptionsPrice(element);
};

// Функция, которая обновляет цену за товары.
// Вызывается при изменении значений или состава товаров
window.updateProductsPrice = () => {
  var productsRetailPrice = 0;
  var productsWholesalePrice = 0;

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
};

// Функция, которая добавляет новый товар в список.
// Клонирует шаблон product-template
window.addProduct = () => {
  const productList = document.getElementById("product-list");
  const productTemplate = document.getElementById("product-template");

  const productClone = productTemplate.content.cloneNode(true);
  productList.append(productClone);

  updateProductsPrice();
};

// Функция, формирующая список товаров для отправки заказа
function buildProductsList() {
  const products = {};
  const productItems = document.getElementsByClassName("product-item");
  [...productItems].forEach((productItem) => {
    const productSelect = productItem.querySelector(".product");
    const selectedProduct = productSelect.options[productSelect.selectedIndex];
    const amount = Number(productItem.querySelector(".amount").value);
    if (products[Number(selectedProduct.value)] === undefined) {
      products[Number(selectedProduct.value)] = 0;
    }
    products[Number(selectedProduct.value)] += amount;
  });
  return products;
}

// Функция, пракильно формирующая описание типа привода и информацию о нем.
function buildDrivePayload(gate) {
  const driveType = getSelectedOption(gate.querySelector(".drive-type")).value;
  if (driveType === "manual") {
    return {
      manual: {
        chainLength: Number(gate.querySelector(".chain-length").value),
      },
    };
  }

  if (driveType === "industrial") {
    return {
      industrial: {
        driveId: Number(gate.querySelector(".drive").value),
      },
    };
  }

  if (driveType === "residential") {
    return {
      residential: {
        driveId: Number(gate.querySelector(".drive").value),
        railId: Number(gate.querySelector(".rail").value),
      },
    };
  }

  return null;
}

// Функция, возвращающа выбранный option у select
function getSelectedOption(select) {
  return select.options[select.selectedIndex];
}

// Обект, хранящий номер типа ворот.
// Нужен при отправке заказа
const gateTypeProto = {
  industrial: 0,
  residential: 1,
};

// Функция для формирования списка ворот с заказе.
// Возвращает массив объектов ворот
function buildGatePayload() {
  const orderGates = [];
  const gates = document.getElementsByClassName("gate-item");
  [...gates].forEach((gate) => {
    // Тип ворот
    const gateType =
      gateTypeProto[getSelectedOption(gate.querySelector(".gate-type")).value];

    // Ширина и высота
    const width = Number(gate.querySelector(".width").value);
    const height = Number(gate.querySelector(".height").value);

    // Высота притолки
    const headroom = Number(gate.querySelector(".headroom").value);

    // Тип подъема
    const liftType = Number(
      getSelectedOption(gate.querySelector(".lift-type")).value,
    );

    // Количество циклов
    const cycleAmount = Number(
      getSelectedOption(gate.querySelector(".cycle-amount")).value,
    );

    // Цвет снаружи
    const colorOut = Number(
      getSelectedOption(gate.querySelector(".color-out")).value,
    );

    // Привод
    const drive = buildDrivePayload(gate);

    // Количество ворот
    const gateAmount = Number(gate.querySelector(".gate-amount").value);

    // Сбор дополнительных опций
    const options = {};
    const optionItems = gate.getElementsByClassName("option-item");
    [...optionItems].forEach((optionItem) => {
      const optionSelect = optionItem.querySelector(".option");
      const selectedOption = optionSelect.options[optionSelect.selectedIndex];
      const amount = Number(optionItem.querySelector(".amount").value);
      if (options[Number(selectedOption.value)] === undefined) {
        options[Number(selectedOption.value)] = 0;
      }
      options[Number(selectedOption.value)] += amount;
    });

    orderGates.push({
      gateType: gateType,
      width: width,
      height: height,
      liftTypeId: liftType,
      colorOutId: colorOut,
      cycleAmountId: cycleAmount,
      options: options,
      headroom: headroom,
      drive: drive,
      amount: gateAmount,
    });
  });

  return orderGates;
}

// Функция дя оформления заказа.
// Собирает все данные с формы, кодирует для protobuf и отправляет на сервер
window.placeOrder = async () => {
  // Формирование payload
  const payload = {
    orderGates: buildGatePayload(),
    products: buildProductsList(),
  };

  // Кодирование в protobuf
  const err = Proto.OrderRequest.verify(payload);
  if (err) throw new Error(err);

  const msg = Proto.OrderRequest.create(payload);
  const buf = Proto.OrderRequest.encode(msg).finish();

  const res = await fetch("/orders", {
    method: "POST",
    headers: { "Content-Type": "application/x-protobuf" },
    body: buf,
  });

  if (!res.ok) throw new Error("Order request failed");

  window.location.href = "/orders";
};

// Функция, удаляющая элемент из списка товаров или дополнительных опций
window.removeItem = (element, itemContainer) => {
  const item = element.closest(itemContainer);
  if (!item) return;

  const additionalOptionsBlock = element.closest(".additional-options");
  const gate = element.closest(".gate-item");

  item.remove();

  if (additionalOptionsBlock) {
    updateOptionsPrice(additionalOptionsBlock);
  } else {
    updateProductsPrice();
  }
};

// Инициализация proto-схем после загрузки
await initProtobuf();
// Чтобы изначально были хотя бы одни ворота
addGateElement();
