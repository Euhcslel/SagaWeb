import {
  gateTypeProto,
  removeItem,
  addGateOption,
  updateOptionsPrice,
  updateGatePrice,
  fmtPrice,
} from "./gate-form-utils.js";
import {
  create,
  toBinary,
  fromBinary,
  GateConfigSchema,
  SizePriceSchema,
} from "./proto_bundle.js";

window.removeItem = removeItem;
window.addGateOption = addGateOption;
window.updateOptionsPrice = updateOptionsPrice;
window.updateGatePrice = updateGatePrice;

// Функция, запрашивающая цену за размер ворот
// принимает параметры: w - ширина, h - высота, t - тип ворот
async function fetchSizePrice(w, h, t) {
  const res = await fetch(`/sizes?width=${w}&height=${h}&gateType=${t}`);

  const buf = await res.arrayBuffer();

  const d = fromBinary(SizePriceSchema, new Uint8Array(buf));

  if (d.price.case === "dealer") {
    return {
      retail: Number(d.price.value.clientPrice) / 100,
      wholesale: Number(d.price.value.dealerPrice) / 100,
    };
  }

  return { retail: Number(d.price.value.clientPrice) / 100, wholesale: 0 };
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

  document.getElementById("size-retail-price").textContent = fmtPrice(price.retail);
  document.getElementById("size-wholesale-price").textContent = price.wholesale
    ? fmtPrice(price.wholesale)
    : "—";

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
  railLabel.hidden = isManual || gateType !== "residential";

  updateGatePrice(document.querySelector(".gate-item"));
};

window.SaveGateConfiguration = async () => {
  let drive;

  switch (document.querySelector(".drive-type").value) {
    case "manual":
      drive = {
        driveType: {
          case: "manual",
          value: { chainLength: Number(document.getElementById("chain-length").value) },
        },
      };
      break;

    case "residential":
      drive = {
        driveType: {
          case: "residential",
          value: {
            driveId: BigInt(document.getElementById("drive").value),
            railId: BigInt(document.getElementById("rail").value),
          },
        },
      };
      break;

    case "industrial":
      drive = {
        driveType: {
          case: "industrial",
          value: { driveId: BigInt(document.getElementById("drive").value) },
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
    optionId: BigInt(optionId),
    amount,
  }));

  const payload = {
    gateType: gateTypeProto[document.getElementById("gate-type").dataset.value],
    width: Number(document.getElementById("width").value),
    height: Number(document.getElementById("height").value),
    liftTypeId: BigInt(document.getElementById("lift-type").value),
    cycleAmountId: BigInt(document.getElementById("cycle-amount").value),
    colorOutId: BigInt(document.getElementById("color-out").value),
    drive: drive,
    options: options,
    headroom: Number(document.getElementById("headroom").value),
    amount: Number(document.getElementById("amount").value),
  };

  const msg = create(GateConfigSchema, payload);
  const buf = toBinary(GateConfigSchema, msg);

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

await updateGateSizePrice();
updateOptionsPrice(document.querySelector(".additional-options"));
updateDriveType();
