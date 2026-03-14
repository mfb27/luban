const els = {
  sessionsList: document.getElementById("sessionsList"),
  modelSelect: document.getElementById("modelSelect"),
  newChatBtn: document.getElementById("newChatBtn"),
  themeToggle: document.getElementById("themeToggle"),
  meBtn: document.getElementById("meBtn"),
  meDialog: document.getElementById("meDialog"),
  meClose: document.getElementById("meClose"),
  meAvatar: document.getElementById("meAvatar"),
  meName: document.getElementById("meName"),
  meAvatar2: document.getElementById("meAvatar2"),
  meName2: document.getElementById("meName2"),
  meId: document.getElementById("meId"),
  authDialog: document.getElementById("authDialog"),
  authClose: document.getElementById("authClose"),
  loginForm: document.getElementById("loginForm"),
  registerForm: document.getElementById("registerForm"),
  loginEmail: document.getElementById("loginEmail"),
  loginPassword: document.getElementById("loginPassword"),
  registerName: document.getElementById("registerName"),
  registerEmail: document.getElementById("registerEmail"),
  registerPassword: document.getElementById("registerPassword"),
  messages: document.getElementById("messages"),
  prompt: document.getElementById("prompt"),
  sendBtn: document.getElementById("sendBtn"),
  errorBar: document.getElementById("errorBar"),
  fileInput: document.getElementById("fileInput"),
  attachments: document.getElementById("attachments"),
  logoutBtn: document.getElementById("logoutBtn"),
};

const state = {
  me: null,
  models: [],
  sessions: [],
  currentSessionId: "",
  currentModelId: "",
  attachmentURLs: [],
  theme: localStorage.getItem("theme") || "light",
  isSending: false,
  isAuthenticated: false,
};

function setError(msg) {
  if (!msg) {
    els.errorBar.classList.add("hidden");
    els.errorBar.textContent = "";
    return;
  }
  els.errorBar.classList.remove("hidden");
  els.errorBar.textContent = msg;
}

function setTheme(theme) {
  state.theme = theme;
  document.documentElement.setAttribute("data-theme", theme);
  localStorage.setItem("theme", theme);
}

async function api(path, options = {}) {
  const token = localStorage.getItem("token");
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  // Add Authorization header if token exists and it's not an auth route
  if (token && !path.startsWith("/api/auth/")) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(path, {
    ...options,
    headers,
  });
  const isJSON = (res.headers.get("content-type") || "").includes("application/json");
  const body = isJSON ? await res.json() : await res.text();

  // Consider 2xx status codes as successful
  if (res.status >= 200 && res.status < 300) {
    return body;
  }

  // For error status codes
  const msg = body?.error || body || `HTTP ${res.status}`;
  throw new Error(msg);
}

async function loadMe() {
  try {
    const me = await api("/api/user/me");
    state.me = me;
    state.isAuthenticated = true;
    // Always use default avatar
    els.meAvatar.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
    els.meName.textContent = me.name;
    els.meAvatar2.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
    els.meName2.textContent = me.name;
    els.meId.textContent = `ID: ${me.id}`;
  } catch (err) {
    // User not authenticated
    state.me = null;
    state.isAuthenticated = false;
    els.meAvatar.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
    els.meName.textContent = "登录/注册";
  }
}

async function loadModels() {
  const models = await api("/api/models");
  state.models = models;
  els.modelSelect.innerHTML = "";
  for (const m of models) {
    const opt = document.createElement("option");
    opt.value = m.id;
    opt.textContent = m.name;
    els.modelSelect.appendChild(opt);
  }
  if (!state.currentModelId && models.length) {
    state.currentModelId = models[0].id;
    els.modelSelect.value = state.currentModelId;
  }
  renderTopbar();
}

function renderTopbar() {
  // 仅保持下拉的选中状态与当前模型同步
  if (els.modelSelect && state.currentModelId) {
    els.modelSelect.value = state.currentModelId;
  }
}

function fmtTime(ts) {
  if (!ts) return "";
  const d = new Date(ts);
  return d.toLocaleString();
}

function renderSessions() {
  els.sessionsList.innerHTML = "";
  for (const s of state.sessions) {
    const div = document.createElement("div");
    div.className = "session-item" + (s.id === state.currentSessionId ? " active" : "");
    div.innerHTML = `
      <div class="session-title"></div>
      <div class="session-meta"></div>
    `;
    div.querySelector(".session-title").textContent = s.title;
    div.querySelector(".session-meta").textContent = fmtTime(s.updated_at);
    div.addEventListener("click", () => selectSession(s.id));
    els.sessionsList.appendChild(div);
  }
}

function renderAttachments() {
  els.attachments.innerHTML = "";
  for (const url of state.attachmentURLs) {
    const span = document.createElement("span");
    span.className = "att";
    span.textContent = url.split("/").slice(-1)[0];
    els.attachments.appendChild(span);
  }
}

// Add a message bubble to the UI
function addMessageBubble(role, content) {
  const wrap = document.createElement("div");
  wrap.className = `msg ${role}`;
  const bubble = document.createElement("div");
  bubble.className = "bubble";
  bubble.textContent = content;
  wrap.appendChild(bubble);
  return wrap;
}

// Add a message without timestamp
function addMessage(role, content) {
  const container = document.createElement("div");
  container.className = `msg ${role}`;

  const bubble = document.createElement("div");
  bubble.className = "bubble";
  bubble.textContent = content;

  container.appendChild(bubble);
  els.messages.appendChild(container);

  // Remove empty class and greeting when messages are added
  els.messages.parentElement.classList.remove("empty");
  const greeting = els.messages.querySelector(".greeting");
  if (greeting) {
    greeting.remove();
  }

  // Smooth scroll to bottom
  els.messages.scrollTo({
    top: els.messages.scrollHeight,
    behavior: 'smooth'
  });
  return container;
}


function clearMessages() {
  els.messages.innerHTML = "";
  els.messages.parentElement.classList.add("empty");

  // Add greeting message when empty
  const greeting = document.createElement("div");
  greeting.className = "greeting";
  greeting.textContent = "你好，我是鲁班";
  els.messages.appendChild(greeting);
}

async function loadSessions() {
  const sessions = await api("/api/sessions");
  state.sessions = sessions;
  renderSessions();
  renderTopbar();
}

async function loadMessages(sessionId) {
  if (!sessionId) return;

  try {
    const msgs = await api(`/api/sessions/${encodeURIComponent(sessionId)}/messages`);
    clearMessages();
    if (msgs.length === 0) {
      // Keep empty state with greeting
      els.messages.parentElement.classList.add("empty");
      const greeting = document.createElement("div");
      greeting.className = "greeting";
      greeting.textContent = "你好，我是鲁班";
      els.messages.appendChild(greeting);
    } else {
      for (const m of msgs) {
        addMessage(m.role, m.content);
      }
    }
  } catch (e) {
    console.error("Failed to load messages:", e);
    // Don't clear messages if loading fails
  }
}

async function selectSession(sessionId) {
  state.currentSessionId = sessionId;
  renderSessions();
  renderTopbar();
  await loadMessages(sessionId);
}

async function newChat() {
  // Check if user is authenticated
  if (!state.isAuthenticated) {
    els.authDialog.showModal();
    return;
  }

  state.currentSessionId = "";
  state.attachmentURLs = [];
  renderAttachments();
  renderSessions();
  renderTopbar();
  clearMessages();
  setError("");
}

async function send() {
  // Check if user is authenticated
  if (!state.isAuthenticated) {
    els.authDialog.showModal();
    return;
  }

  if (state.isSending) return;
  setError("");
  const content = els.prompt.value.trim();
  if (!content) return;

  const modelId = els.modelSelect.value || state.currentModelId;
  state.currentModelId = modelId;

  state.isSending = true;
  els.sendBtn.disabled = true;

  // Disable form inputs during sending
  els.prompt.disabled = true;
  els.fileInput.disabled = true;

  try {
    // Add user message
    addMessage("user", content);
    els.prompt.value = "";

    // Create assistant message container
    const assistantContainer = document.createElement("div");
    assistantContainer.className = "msg assistant";

    const bubble = document.createElement("div");
    bubble.className = "bubble";
    bubble.textContent = "";

    assistantContainer.appendChild(bubble);
    els.messages.appendChild(assistantContainer);

    // Prepare request body
    const reqBody = JSON.stringify({
      session_id: state.currentSessionId,
      content,
      model_id: modelId,
      attachment_urls: state.attachmentURLs,
    });

    // Get token for authenticated request
    const token = localStorage.getItem("token");
    const headers = { "Content-Type": "application/json" };

    // Add Authorization header if token exists
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    // Call streaming API
    const response = await fetch("/api/chat", {
      method: "POST",
      headers,
      body: reqBody,
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || `HTTP ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = "";

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split("\n");
      buffer = lines.pop() || ""; // Keep incomplete line in buffer

      for (const line of lines) {
        if (!line.trim() || !line.startsWith("data:")) continue;

        const data = line.slice(5).trim();
        if (!data) continue;

        try {
          const event = JSON.parse(data);

          switch (event.type) {
            case "session":
              // Update session ID only if it's not already set
              if (!state.currentSessionId) {
                state.currentSessionId = event.session_id;
              }
              break;

            case "content":
              // Append content to assistant bubble
              bubble.textContent += event.data;
              // Smooth scroll to bottom
              requestAnimationFrame(() => {
                els.messages.scrollTo({
                  top: els.messages.scrollHeight,
                  behavior: 'smooth'
                });
              });
              break;

            case "done":
              // Update session list
              await loadSessions();
              break;

            case "error":
              throw new Error(event.error || "Unknown error");
          }
        } catch (e) {
          if (e instanceof SyntaxError) {
            // Skip invalid JSON
            continue;
          }
          throw e;
        }
      }
    }

    // Clear attachments after successful send
    state.attachmentURLs = [];
    renderAttachments();
  } catch (e) {
    setError(e.message || String(e));
    // Keep all messages even on error
  } finally {
    state.isSending = false;
    els.sendBtn.disabled = false;
    els.prompt.disabled = false;
    els.fileInput.disabled = false;
    renderTopbar();
  }
}

async function uploadFile(file) {
  setError("");
  const fd = new FormData();
  fd.append("file", file);

  // Get token for authenticated request
  const token = localStorage.getItem("token");
  const headers = {};

  // Add Authorization header if token exists
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch("/api/upload", {
    method: "POST",
    headers,
    body: fd
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error || `HTTP ${res.status}`);
  state.attachmentURLs.push(body.url);
  renderAttachments();
}

function bindEvents() {
  els.newChatBtn.addEventListener("click", newChat);

  els.modelSelect.addEventListener("change", () => {
    state.currentModelId = els.modelSelect.value;
    renderTopbar();
  });

  els.themeToggle.addEventListener("click", () => {
    setTheme(state.theme === "dark" ? "light" : "dark");
  });

  // Auth events
  els.meBtn.addEventListener("click", () => {
    if (state.isAuthenticated) {
      els.meDialog.showModal();
    } else {
      els.authDialog.showModal();
      setTimeout(() => {
        els.loginEmail.focus();
      }, 100);
    }
  });
  els.authClose.addEventListener("click", () => els.authDialog.close());
  els.meClose.addEventListener("click", () => els.meDialog.close());

  els.loginForm.addEventListener("submit", handleLogin);
  els.registerForm.addEventListener("submit", handleRegister);

  // Tab switching
  document.querySelectorAll(".auth-tab").forEach(tab => {
    tab.addEventListener("click", () => switchAuthTab(tab.dataset.tab));
  });

  // Password strength indicator
  els.registerPassword.addEventListener("input", (e) => {
    checkPasswordStrength(e.target.value);
  });

  els.sendBtn.addEventListener("click", send);
  els.prompt.addEventListener("keydown", (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      send();
    }
  });

  // Update send button state when prompt content changes
  els.prompt.addEventListener("input", () => {
    els.sendBtn.disabled = els.prompt.value.trim() === '' || state.isSending;
  });

  els.fileInput.addEventListener("change", async () => {
    const file = els.fileInput.files?.[0];
    if (!file) return;
    els.fileInput.value = "";
    try {
      await uploadFile(file);
    } catch (e) {
      setError(e.message || String(e));
    }
  });

  // Logout button
  els.logoutBtn.addEventListener("click", logout);
}

async function boot() {
  setTheme(state.theme);
  bindEvents();

  // Initialize empty state
  els.messages.parentElement.classList.add("empty");

  // Add greeting message on load
  const greeting = document.createElement("div");
  greeting.className = "greeting";
  greeting.textContent = "你好，我是鲁班";
  els.messages.appendChild(greeting);

  try {
    // Load models (no auth required)
    await loadModels();

    // Check authentication and load user data
    await checkAuth();

    if (state.isAuthenticated && state.me && state.me.name) {
      // Update greeting with user name
      greeting.textContent = `你好，${state.me.name}`;
    }

    // Load sessions if authenticated
    if (state.isAuthenticated) {
      await loadSessions();
    }

    renderAttachments();
  } catch (e) {
    setError(e.message || String(e));
  }
}

// Auth functions
function handleLogin(e) {
  e.preventDefault();

  const email = els.loginEmail.value;
  const password = els.loginPassword.value;

  api("/api/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  })
    .then(res => {
      // Login returns 200 with response containing token and user info
      state.me = {
        id: res.user_id,
        name: res.name,
        email: res.email
      };
      state.isAuthenticated = true;
      localStorage.setItem("token", res.token);
      els.authDialog.close();
      // Update avatar to default
      els.meAvatar.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
      updateUIAfterAuth();
    })
    .catch(err => {
      setError(err.message || "登录失败");
    });
}

function checkPasswordStrength(password) {
  const strengthBar = document.querySelector('.strength-bar');
  const strengthText = document.querySelector('.strength-text');

  if (!password) {
    strengthBar.className = 'strength-bar';
    strengthText.className = 'strength-text';
    strengthText.textContent = '';
    return;
  }

  let strength = 0;
  const checks = [
    password.length >= 6,
    password.length >= 8,
    /[a-z]/.test(password),
    /[A-Z]/.test(password),
    /[0-9]/.test(password),
    /[^a-zA-Z0-9]/.test(password)
  ];

  strength = checks.filter(Boolean).length;

  if (strength <= 2) {
    strengthBar.className = 'strength-bar weak';
    strengthText.className = 'strength-text weak';
    strengthText.textContent = '弱';
  } else if (strength <= 4) {
    strengthBar.className = 'strength-bar medium';
    strengthText.className = 'strength-text medium';
    strengthText.textContent = '中';
  } else {
    strengthBar.className = 'strength-bar strong';
    strengthText.className = 'strength-text strong';
    strengthText.textContent = '强';
  }
}

function handleRegister(e) {
  e.preventDefault();

  const name = els.registerName.value;
  const email = els.registerEmail.value;
  const password = els.registerPassword.value;

  api("/api/auth/register", {
    method: "POST",
    body: JSON.stringify({ name, email, password }),
  })
    .then(res => {
      // Registration returns 201 with response containing token and user info
      state.me = {
        id: res.user_id,
        name: res.name,
        email: res.email
      };
      state.isAuthenticated = true;
      localStorage.setItem("token", res.token);
      els.authDialog.close();
      // Update avatar to default
      els.meAvatar.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
      updateUIAfterAuth();
    })
    .catch(err => {
      setError(err.message || "注册失败");
    });
}

function switchAuthTab(tab) {
  document.querySelectorAll(".auth-tab").forEach(t => {
    t.classList.toggle("active", t.dataset.tab === tab);
  });

  els.loginForm.classList.toggle("hidden", tab !== "login");
  els.registerForm.classList.toggle("hidden", tab !== "register");

  // Auto-focus appropriate input
  setTimeout(() => {
    if (tab === "login") {
      els.loginEmail.focus();
    } else {
      els.registerName.focus();
    }
  }, 100);
}

function updateUIAfterAuth() {
  // Update user info in UI
  els.meName.textContent = state.me.name;
  // Always use default avatar for now
  els.meAvatar.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";

  // Update sidebar
  document.querySelector('.me-sub').textContent = '个人信息';

  // Load user data
  fetchUserProfile();
  loadSessions();
}

function fetchUserProfile() {
  api("/api/user/me", {
    headers: {
      "Authorization": `Bearer ${localStorage.getItem("token")}`,
    },
  })
    .then(res => {
      if (res.id) {
        state.me = res;
        els.meName.textContent = res.name;
        els.meAvatar.src = res.avatar_url || "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
        els.meName2.textContent = res.name;
        els.meAvatar2.src = res.avatar_url || "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
        els.meId.textContent = `ID: ${res.id}`;
      }
    })
    .catch(err => {
      console.error("Failed to fetch user profile:", err);
    });
}

async function checkAuth() {
  const token = localStorage.getItem("token");
  if (token) {
    try {
      // Verify token is valid
      const res = await api("/api/user/me", {
        headers: {
          "Authorization": `Bearer ${token}`,
        },
      });

      if (res.id) {
        state.me = res;
        state.isAuthenticated = true;
        updateUIAfterAuth();
      } else {
        localStorage.removeItem("token");
        state.isAuthenticated = false;
        els.meName.textContent = "登录/注册";
        els.meAvatar.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
      }
    } catch (error) {
      localStorage.removeItem("token");
      state.isAuthenticated = false;
      els.meName.textContent = "登录/注册";
      els.meAvatar.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
    }
  } else {
    // No token, not authenticated
    state.isAuthenticated = false;
    els.meName.textContent = "登录/注册";
    // Avatar already has default src from HTML
  }
}

function logout() {
  // Clear authentication data
  localStorage.removeItem("token");
  state.isAuthenticated = false;
  state.me = null;
  state.sessions = [];
  state.currentSessionId = "";

  // Reset UI
  els.meDialog.close();
  els.meName.textContent = "登录/注册";
  els.meAvatar.src = "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png";
  els.sessionsList.innerHTML = "";

  // Clear messages and show greeting
  clearMessages();

  // Update sidebar
  document.querySelector('.me-sub').textContent = '个人中心';

  // Reset form
  els.prompt.value = "";
  els.prompt.disabled = false;
  els.sendBtn.disabled = true;
}

boot();
