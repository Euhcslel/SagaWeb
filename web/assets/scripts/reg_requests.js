async function confirmRegRequest(requestId) {
  const result = await Swal.fire({
    title: "Подтвердить заявку?",
    text: "Вы уверены, что хотите подтвердить заявку?",
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Да, подтвердить",
    cancelButtonText: "Отмена",
  });

  if (result.isConfirmed) {
    try {
      const response = await fetch(`/dealers/requests/${requestId}/confirm`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Ошибка подтверждения заявки");
      }

      window.location.reload();
    } catch (error) {
      Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
    }
  }
}

async function rejectRegRequest(requestId) {
  const result = await Swal.fire({
    title: "Отказать в заявке?",
    text: "Вы уверены, что хотите отклонить заявку? Это действие необратимо.",
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "Да, отказать",
    cancelButtonText: "Отмена",
  });

  if (result.isConfirmed) {
    try {
      const response = await fetch(`/dealers/requests/${requestId}/reject`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Ошибка отклонения заявки");
      }

      window.location.reload();
    } catch (error) {
      Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
    }
  }
}
