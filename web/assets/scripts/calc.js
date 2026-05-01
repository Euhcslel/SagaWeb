import {
  gateTypeProto,
  removeItem,
  addGateOption,
  addProduct,
  updateOptionsPrice,
  updateProductsPrice,
  updateGatePrice as sharedUpdateGatePrice,
} from "./gate-form-utils.js";

window.removeItem = removeItem;
window.addGateOption = addGateOption;
window.updateOptionsPrice = updateOptionsPrice;
window.addProduct = addProduct;
window.updateProductsPrice = updateProductsPrice;
window.updateGatePrice = (gate) => {
  sharedUpdateGatePrice(gate);
  updateOrderPrice();
}

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
    (gate?.dataset.retailPrice || 0).toFixed(2);

  document.getElementById("gate-wholesale-price").textContent =
    (gate?.dataset.wholesalePrice || 0).toFixed(2);
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
window.updateGateType = async (select) => {
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
  await updateGateSizePrice(width);
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
  orderRetailPriceElement.textContent = totalRetailPrice.toFixed(2);

  const orderWholesalePriceElement = document.getElementById(
    "order-wholesale-price",
  );
  orderWholesalePriceElement.textContent = totalWholesalePrice.toFixed(2);
};

// Функция, формирующая список товаров для отправки заказа
function buildProductsList() {
  const productsMap = {};

  const productItems = document.getElementsByClassName("product-item");

  [...productItems].forEach((productItem) => {
    const productSelect = productItem.querySelector(".product");
    const selectedProduct = productSelect.options[productSelect.selectedIndex];

    const productId = Number(selectedProduct.value);
    const amount = Number(productItem.querySelector(".amount").value);

    if (productsMap[productId] === undefined) {
      productsMap[productId] = 0;
    }

    productsMap[productId] += amount;
  });

  return Object.entries(productsMap).map(([productId, amount]) => ({
    productId: Number(productId),
    amount: amount,
  }));
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
    const optionsMap = {};

    const optionItems = gate.getElementsByClassName("option-item");

    [...optionItems].forEach((optionItem) => {
      const optionSelect = optionItem.querySelector(".option");
      const selectedOption = optionSelect.options[optionSelect.selectedIndex];
      const optionId = Number(selectedOption.value);
      const amount = Number(optionItem.querySelector(".amount").value);

      if (optionsMap[optionId] === undefined) {
        optionsMap[optionId] = 0;
      }

      optionsMap[optionId] += amount;
    });

    const options = Object.entries(optionsMap).map(([optionId, amount]) => ({
      optionId: Number(optionId),
      amount: amount,
    }));

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
  const payload = {
    orderGates: buildGatePayload(),
    products: buildProductsList(),
  };

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

// Инициализация proto-схем после загрузки
await initProtobuf();
// Чтобы изначально были хотя бы одни ворота
addGateElement();
