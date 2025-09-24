const API_BASE = "/api/notify";

async function parseJSON(resp) {
	const payload = await resp.json().catch(() => null);
	if (!payload) return null;
	if (payload.data !== undefined) return payload.data;
	if (payload.result !== undefined) return payload.result;
	return payload;
}

function renderEmptyState(container) {
	container.innerHTML = `<div class="empty">Пока нет уведомлений</div>`;
}

function formatDate(value) {
	return value || "—";
}

async function loadNotifications() {
	const listDiv = document.getElementById("list");
	listDiv.innerHTML = "";

	const resp = await fetch(API_BASE + "/");
	const data = await parseJSON(resp);
	if (!data || data.length === 0) {
		return renderEmptyState(listDiv);
	}

	data.forEach((n) => {
		const div = document.createElement("div");
		div.className = "notif";
		div.innerHTML = `
			<div>
				<p><b>ID:</b> ${n.id}</p>
				<p><b>Message:</b> ${n.message}</p>
				<p><b>Channel:</b> <span class="badge">${n.channel}</span></p>
				<p><b>To:</b> ${n.to}</p>
				<p><b>SendAt:</b> ${formatDate(n.send_at)}</p>
				<p><b>Status:</b> <span class="badge">${n.status}</span></p>
			</div>
			<div class="row-end">
				<button class="btn btn-danger" data-id="${n.id}">Удалить</button>
			</div>
		`;
		div.querySelector("button").addEventListener("click", () => deleteNotif(n.id));
		listDiv.appendChild(div);
	});
}

async function deleteNotif(id) {
	await fetch(API_BASE + "/" + id, { method: "DELETE" });
	await loadNotifications();
}

function attachFormHandler() {
	const form = document.getElementById("create-form");
	form.addEventListener("submit", async (e) => {
		e.preventDefault();
		const notif = {
			message: document.getElementById("message").value,
			channel: document.getElementById("channel").value,
			to: document.getElementById("to").value,
			send_at: document.getElementById("sendat").value,
			retries: parseInt(document.getElementById("retries").value)
		};
		await fetch(API_BASE + "/", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(notif)
		});
		form.reset();
		await loadNotifications();
	});
}

window.addEventListener("DOMContentLoaded", () => {
	attachFormHandler();
	loadNotifications();
});


