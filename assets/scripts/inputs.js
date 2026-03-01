// Custom SELECT
document.addEventListener("mousedown", (e) => {
  const option = e.target.closest(".select ul button");
  if (!option) return;

  e.preventDefault();

  const select = option.closest(".select");
  const trigger = select.querySelector("button");
  const label = trigger.querySelector("span");
  const input = select.querySelector("input[type='hidden']");

  label.textContent = option.textContent;

  [...input.attributes].forEach((attr) => {
    if (attr.name.startsWith("data-")) {
      input.removeAttribute(attr.name);
    }
  });

  [...option.attributes].forEach((attr) => {
    if (attr.name.startsWith("data-")) {
      input.setAttribute(attr.name, attr.value);
    }
  });

  input.value = option.dataset.id ?? "";

  trigger.blur();

  updateTotalPrice();
});

// Custom PASSWORD input and DATE input and NUMBER input
document.addEventListener("click", (e) => {
  if (e.target.closest(".password button")) {
    const button = e.target.closest(".password button");
    if (!button) return;

    const wrapper = button.closest(".password");
    const input = wrapper.querySelector("input");
    const img = button.querySelector("img");

    const visible = input.type === "text";
    input.type = visible ? "password" : "text";

    img.src = visible
      ? "/assets/images/eye_open.svg"
      : "/assets/images/eye_closed.svg";
  } else if (e.target.closest(".date button")) {
    const button = e.target.closest(".date button");
    if (!button) return;

    const input = button.closest(".date").querySelector("input[type='date']");

    if (input?.showPicker) {
      input.showPicker();
    } else {
      input.focus();
    }
  } else if (e.target.closest(".number button")) {
    const button = e.target.closest(".number button");
    if (!button) return;

    const wrapper = button.closest(".number");
    const input = wrapper.querySelector("input");

    button.dataset.step === "up" ? input.stepUp() : input.stepDown();
  }
});

// Custom PHONE input
document.addEventListener("input", (e) => {
  const input = e.target.closest(".phone input");
  if (!input) return;

  let digits = input.value.replace(/\D/g, "").slice(0, 10);

  let result = "";
  if (digits.length > 0) result += "(" + digits.slice(0, 3);
  if (digits.length >= 4) result += ") " + digits.slice(3, 6);
  if (digits.length >= 7) result += "−" + digits.slice(6, 8);
  if (digits.length >= 9) result += "−" + digits.slice(8, 10);

  input.value = result;
});
