// Функция, которая обновляет цену на текущие ворота.
// Вызывается при изменении занчений select и input, которые влияют на стоимость ворот
window.updateGatePrice = () => {
  var retailGatePrice = 0;
  var wholesaleGatePrice = 0;

  const gate = document.querySelector(".configuration");

  const size = gate.querySelector(".sizes");
  const sizeRetailPrice = Number(size.dataset.retailPrice || 0);
  const sizeWholesalePrice = Number(size.dataset.wholesalePrice || 0);

  const liftType = gate.querySelector("#lift-type");
  const selectedLiftType = liftType.options[liftType.selectedIndex];
  const liftTypeRetailMarkup = Number(selectedLiftType.dataset.retailMarkup);
  const liftTypeWholesaleMarkup = Number(
    selectedLiftType.dataset.wholesaleMarkup,
  );

  const cycleAmount = gate.querySelector("#cycle-amount");
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
      const drive = gate.querySelector("#drive");
      const selectedDrive = drive.options[drive.selectedIndex];
      const rail = gate.querySelector("#rail");
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
      const drive = gate.querySelector("#drive");
      const selectedDrive = drive.options[drive.selectedIndex];
      driveRetailPrice += Number(selectedDrive.dataset.retailPrice);
      driveWholesalePrice += Number(selectedDrive.dataset.wholesalePrice);
      break;
    }
  }

  var optionList = document.querySelector(".option-list");
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
};

// Proto-схемы
let Proto = {
  SizePrice: null,
  GateConfig: null,
};

// Функция, инициализирующая proto-схемы
async function initProtobuf() {
  const pricesRoot = await protobuf.load("/api/proto/prices.proto");
  Proto.SizePrice = pricesRoot.lookupType("proto.SizePrice");

  const gateRoot = await protobuf.load("/api/proto/order.proto");
  Proto.GateConfig = gateRoot.lookupType("proto.GateConfig");
}

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

// Функция, срабатывающая при изменении значений ширины и высоты
// запрашивает цены клиента и дилера у сервера и присваивает их блоку sizes
window.updateGateSizePrice = async () => {
  const width = document.getElementById("width");
  const height = document.getElementById("height");
  const gateType = document.getElementById("gate-type");

  const sizes = document.querySelector(".sizes");

  const price = await fetchSizePrice(
    width.value,
    height.value,
    gateType.dataset.value,
  );

  sizes.dataset.retailPrice = price.retail;
  sizes.dataset.wholesalePrice = price.wholesale;

  updateGatePrice();
};

// Функция, которая обновляет цену за дополнительные опции у ворот.
// Вызывается при изменении значений или состава дополнительных опций
window.updateOptionsPrice = () => {
  var optionsRetailPrice = 0;
  var optionsWholesalePrice = 0;

  const optionList = document.querySelector(".option-list");
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

  updateGatePrice();
};

// Функция, срабатывающая при изменении привода
// показывает нужные поля ввода и скрывает ненужные
window.updateDriveType = () => {
  const select = document.querySelector(".drive-type");
  const isManual = select.value === "manual";

  const chainLabel = document.getElementById("chain-length").closest("label");
  const driveAutoBlock = document.querySelector(".drive-auto");
  const railLabel = document.getElementById("rail").closest("label");
  const gateType = document.getElementById("gate-type").dataset.value;

  chainLabel.hidden = !isManual;
  driveAutoBlock.hidden = isManual;
  railLabel.hidden = isManual || gateType !== "res";

  updateGatePrice();
};

// Обект, хранящий номер типа ворот.
// Нужен при отправке заказа
const gateTypeProto = {
  industrial: 0,
  residential: 1,
};

window.SaveGateConfiguration = async () => {
  // Формирование payload
  let drive;
  switch (document.querySelector(".drive-type").value) {
    case "manual":
      drive = {
        manual: {
          chainLength: Number(document.getElementById("chain-length").value),
        },
      };
      break;
    case "residential":
      drive = {
        residential: {
          driveId: Number(document.getElementById("drive").value),
          railId: Number(document.getElementById("rail").value),
        },
      };
      break;
    case "industrial":
      drive = {
        industrial: {
          driveId: Number(document.getElementById("drive").value),
        },
      };
      break;
  }

  const payload = {
    gateType: gateTypeProto[document.getElementById("gate-type").dataset.value],
    width: Number(document.getElementById("width").value),
    height: Number(document.getElementById("height").value),
    liftTypeId: Number(document.getElementById("lift-type").value),
    cycleAmountId: Number(document.getElementById("cycle-amount").value),
    colorOutId: Number(document.getElementById("color-out").value),
    drive: drive,
    options: [...document.querySelectorAll(".option-item")].map(
      (optionItem) => ({
        optionId: Number(optionItem.querySelector(".option").value),
        amount: Number(optionItem.querySelector(".amount").value),
      }),
    ),
    headroom: Number(document.getElementById("headroom").value),
    amount: Number(document.getElementById("amount").value),
  };
  
  // Кодирование в protobuf
  const err = Proto.GateConfig.verify(payload);
  if (err) throw new Error(err);

  const msg = Proto.GateConfig.create(payload);
  const buf = Proto.GateConfig.encode(msg).finish();

  const path = window.location.pathname;

  const res = await fetch(path, {
    method: "PUT",
    headers: { "Content-Type": "application/x-protobuf" },
    body: buf,
  });

  if (!res.ok) throw new Error("Order request failed");

  window.location.reload();
};

await initProtobuf();
await updateGateSizePrice();
updateOptionsPrice();
updateDriveType();
