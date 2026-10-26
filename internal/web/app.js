const apiBase = "/api/v1";
const tokenKey = "token";

function makeLink(label, href) {
    const link = document.createElement("a");
    link.href = href;
    link.textContent = label;
    return link;
}

function renderNavbar(page) {
    const navbar = document.querySelector("#navbar");
    const brand = makeLink("Short URLs", page === "dashboard" ? "/dashboard" : "/");
    brand.className = "brand";
    navbar.append(brand);

    const actions = document.createElement("div");
    actions.className = "nav-actions";
    actions.append(makeLink("API docs", "/swagger/"));
    const github = makeLink("GitHub", "https://github.com/isa0-gh/short-urls");
    github.target = "_blank";
    github.rel = "noopener noreferrer";
    actions.append(github);

    if (page === "dashboard") {
        const logout = document.createElement("button");
        logout.type = "button";
        logout.textContent = "Log out";
        logout.addEventListener("click", () => {
            localStorage.removeItem(tokenKey);
            window.location.replace("/login");
        });
        actions.append(logout);
    } else if (page === "register") {
        actions.append(makeLink("Log in", "/login"));
    } else {
        actions.append(makeLink("Create account", "/register"));
    }

    navbar.append(actions);
}

function setMessage(element, text, success = false) {
    element.textContent = text;
    element.classList.toggle("success", success);
}

async function request(path, options = {}, authenticated = true) {
    const headers = new Headers(options.headers || {});
    const token = authenticated ? localStorage.getItem(tokenKey) : null;
    if (token) headers.set("Authorization", `Bearer ${token}`);
    if (options.body) headers.set("Content-Type", "application/json");

    const response = await fetch(`${apiBase}${path}`, { ...options, headers });
    const text = await response.text();
    let result = null;
    if (text) {
        try {
            result = JSON.parse(text);
        } catch {
            result = null;
        }
    }

    if (response.status === 401 && token) {
        localStorage.removeItem(tokenKey);
        window.location.replace("/login");
    }
    if (!response.ok) {
        throw new Error(result?.error || `Request failed (${response.status}).`);
    }
    return result;
}

function valuesFrom(form) {
    return Object.fromEntries(new FormData(form).entries());
}

function setBusy(form, busy, label) {
    const button = form.querySelector("button[type='submit']");
    button.disabled = busy;
    if (busy) {
        button.dataset.label = button.textContent;
        button.textContent = label;
    } else if (button.dataset.label) {
        button.textContent = button.dataset.label;
        delete button.dataset.label;
    }
}

function setupLogin() {
    const form = document.querySelector("#login-form");
    const message = document.querySelector("#message");
    if (new URLSearchParams(window.location.search).has("registered")) {
        setMessage(message, "Account created. You can log in now.", true);
    }

    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        setMessage(message, "");
        setBusy(form, true, "Logging in...");
        try {
            const result = await request("/login", {
                method: "POST",
                body: JSON.stringify(valuesFrom(form)),
            }, false);
            localStorage.setItem(tokenKey, result.token);
            window.location.assign("/dashboard");
        } catch (error) {
            setMessage(message, error.message || "Could not connect. Please try again.");
        } finally {
            setBusy(form, false);
        }
    });
}

function setupRegister() {
    const form = document.querySelector("#register-form");
    const message = document.querySelector("#message");
    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        setMessage(message, "");
        setBusy(form, true, "Creating account...");
        try {
            await request("/register", {
                method: "POST",
                body: JSON.stringify(valuesFrom(form)),
            }, false);
            window.location.assign("/login?registered=1");
        } catch (error) {
            setMessage(message, error.message || "Could not create account. Please try again.");
        } finally {
            setBusy(form, false);
        }
    });
}

function shortURL(shortCode) {
    return new URL(`/s/${encodeURIComponent(shortCode)}`, window.location.origin).href;
}

function isWebURL(value) {
    try {
        return ["http:", "https:"].includes(new URL(value).protocol);
    } catch {
        return false;
    }
}

function makeCell(content) {
    const cell = document.createElement("td");
    if (content instanceof Node) cell.append(content);
    else cell.textContent = content;
    return cell;
}

function renderLinks(links) {
    const table = document.querySelector("#links-table");
    table.replaceChildren();
    if (!links.length) {
        const row = document.createElement("tr");
        const cell = makeCell("No links yet. Create one above.");
        cell.colSpan = 4;
        cell.className = "empty-cell";
        row.append(cell);
        table.append(row);
        return;
    }

    for (const link of links) {
        const row = document.createElement("tr");
        const publicURL = shortURL(link.short_code);
        const shortAnchor = makeLink(publicURL, publicURL);
        shortAnchor.target = "_blank";
        shortAnchor.rel = "noopener noreferrer";
        row.append(makeCell(shortAnchor));

        const destination = isWebURL(link.destination_url)
            ? makeLink(link.destination_url, link.destination_url)
            : document.createTextNode(link.destination_url);
        if (destination instanceof HTMLAnchorElement) {
            destination.target = "_blank";
            destination.rel = "noopener noreferrer";
        }
        row.append(makeCell(destination));
        row.append(makeCell(new Date(link.created_at).toLocaleDateString()));

        const actions = document.createElement("div");
        actions.className = "row-actions";

        const copy = document.createElement("button");
        copy.type = "button";
        copy.textContent = "Copy";
        copy.addEventListener("click", async () => {
            try {
                await navigator.clipboard.writeText(publicURL);
                setMessage(document.querySelector("#links-message"), "Short link copied.", true);
            } catch {
                setMessage(document.querySelector("#links-message"), "Could not copy link.");
            }
        });
        actions.append(copy);

        const edit = document.createElement("button");
        edit.type = "button";
        edit.textContent = "Edit";
        edit.addEventListener("click", async () => {
            const value = window.prompt("New destination URL", link.destination_url);
            if (value === null || !value.trim()) return;
            const destinationURL = value.trim();
            if (!isWebURL(destinationURL)) {
                setMessage(document.querySelector("#links-message"), "Enter a valid http or https URL.");
                return;
            }
            try {
                await request(`/link/${encodeURIComponent(link.id)}`, {
                    method: "POST",
                    body: JSON.stringify({ destination_url: destinationURL }),
                });
                await loadDashboard();
            } catch (error) {
                setMessage(document.querySelector("#links-message"), error.message);
            }
        });
        actions.append(edit);

        const remove = document.createElement("button");
        remove.type = "button";
        remove.textContent = "Delete";
        remove.className = "delete-button";
        remove.addEventListener("click", async () => {
            if (!window.confirm(`Delete ${link.short_code}?`)) return;
            try {
                await request(`/link/${encodeURIComponent(link.id)}`, { method: "DELETE" });
                await loadDashboard();
            } catch (error) {
                setMessage(document.querySelector("#links-message"), error.message);
            }
        });
        actions.append(remove);
        row.append(makeCell(actions));
        table.append(row);
    }
}

async function loadDashboard() {
    const message = document.querySelector("#links-message");
    setMessage(message, "");
    try {
        const [profile, links] = await Promise.all([
            request("/profile"),
            request("/links"),
        ]);
        document.querySelector("#profile").textContent = `${profile.username} · ${profile.email}`;
        renderLinks(Array.isArray(links) ? links : []);
    } catch (error) {
        setMessage(message, error.message || "Could not load your links.");
    }
}

function setupDashboard() {
    if (!localStorage.getItem(tokenKey)) {
        window.location.replace("/login");
        return;
    }

    const form = document.querySelector("#create-form");
    const message = document.querySelector("#create-message");
    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        setMessage(message, "");
        const values = valuesFrom(form);
        if (!isWebURL(values.destination_url)) {
            setMessage(message, "Enter a valid http or https URL.");
            return;
        }
        setBusy(form, true, "Creating...");
        try {
            await request("/create", {
                method: "POST",
                body: JSON.stringify(values),
            });
            form.reset();
            setMessage(message, "Link created.", true);
            await loadDashboard();
        } catch (error) {
            setMessage(message, error.message || "Could not create link.");
        } finally {
            setBusy(form, false);
        }
    });
    loadDashboard();
}

const page = document.body.dataset.page;
renderNavbar(page);

if (page === "login") setupLogin();
if (page === "register") setupRegister();
if (page === "dashboard") setupDashboard();