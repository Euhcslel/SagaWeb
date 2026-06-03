import {
  gateTypeProto,
  removeItem,
  addGateOption,
  updateOptionsPrice,
  updateGatePrice,
  fmtPrice,
} from "./gate-form-utils.js";

window.removeItem = removeItem;
window.addGateOption = addGateOption;
window.updateOptionsPrice = updateOptionsPrice;
window.updateGatePrice = updateGatePrice;

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

  updateGatePrice(document.querySelector(".gate-item"));
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

  updateGatePrice(document.querySelector(".gate-item"));
};

window.SaveGateConfiguration = async () => {
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

  const optionMap = new Map();

  document.querySelectorAll(".option-item").forEach((optionItem) => {
    const optionId = Number(optionItem.querySelector(".option").value);
    const amount = Number(optionItem.querySelector(".amount").value);

    if (!optionId || Number.isNaN(amount) || amount <= 0) {
      return;
    }

    const currentAmount = optionMap.get(optionId) || 0;
    optionMap.set(optionId, currentAmount + amount);
  });

  const options = Array.from(optionMap.entries()).map(([optionId, amount]) => ({
    optionId,
    amount,
  }));

  const payload = {
    gateType: gateTypeProto[document.getElementById("gate-type").dataset.value],
    width: Number(document.getElementById("width").value),
    height: Number(document.getElementById("height").value),
    liftTypeId: Number(document.getElementById("lift-type").value),
    cycleAmountId: Number(document.getElementById("cycle-amount").value),
    colorOutId: Number(document.getElementById("color-out").value),
    drive: drive,
    options: options,
    headroom: Number(document.getElementById("headroom").value),
    amount: Number(document.getElementById("amount").value),
  };

  const err = Proto.GateConfig.verify(payload);
  if (err) {
    throw new Error(err);
  }

  const msg = Proto.GateConfig.create(payload);
  const buf = Proto.GateConfig.encode(msg).finish();

  const path = window.location.pathname;

  const res = await fetch(path, {
    method: "PUT",
    headers: {
      "Content-Type": "application/x-protobuf",
    },
    body: buf,
  });

  if (!res.ok) {
    throw new Error("Order request failed");
  }

  window.location.reload();
};

await initProtobuf();
await updateGateSizePrice();
updateOptionsPrice(document.querySelector(".additional-options"));
updateDriveType();
