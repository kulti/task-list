export function showSuccessAlert(text: string): void {
  showAlert("success", text);
}

export function showErrorAlert(text: string): void {
  showAlert("danger", text);
}

export function showErrorAlertWithRefresh(text: string): void {
  showAlert("danger", text, 5);
}

function showAlert(type: string, text: string, refreshSec?: number): void {
  const alertCloseBtn = document.createElement("button");
  alertCloseBtn.type = "button";
  alertCloseBtn.className = "close";
  alertCloseBtn.setAttribute("data-dismiss", "alert");
  alertCloseBtn.setAttribute("aria-label", "Close");
  alertCloseBtn.innerHTML = '<span aria-hidden="true">&times;</span>';

  const alertDiv = document.createElement("div");
  alertDiv.className = "alert alert-dismissible fade show alert-" + type;
  alertDiv.setAttribute("role", "alert");
  alertDiv.innerText = text;
  alertDiv.appendChild(alertCloseBtn);

  $("#alerts").append(alertDiv);

  if (refreshSec) {
    alertDiv.innerText =
      text + " - refresh in " + refreshSec.toString() + "seconds";
    let i = refreshSec - 1;
    setInterval(() => {
      alertDiv.innerText = text + " - refresh in " + i.toString() + " seconds";
      i--;
      if (i === 0) {
        window.location.reload();
      }
    }, 1000);
  } else {
    setTimeout(() => {
      alertDiv.remove();
    }, 2000);
  }
}
