
(function() {
  // Authentication state
  var authState = {
    token: null,
    user: null,
    expiresAt: null,
    isAuthenticated: false
  };

  window.addEventListener('DOMContentLoaded', function() {
    initControls();
    window.setInterval(updateTimers, 1000);
    initAuth();
  });

  var tooltipDict = {};
  var tooltipIdx = 1;
  var spamoor = window.spamoor = {
    initControls: initControls,
    renderRecentTime: renderRecentTime,
    tooltipDict: tooltipDict,
    getAuthToken: getAuthToken,
    isAuthenticated: isAuthenticated,
    authFetch: authFetch,
    authState: authState,
  };

  var refreshTimer = null;

  // Initialize authentication
  function initAuth() {
    /* Check if auth is disabled via server config */
    if (window.spamoorConfig && window.spamoorConfig.authDisabled) {
      authState.token = null;
      authState.user = null;
      authState.expiresAt = null;
      authState.isAuthenticated = true;
      authState.authDisabled = true;
      sessionStorage.removeItem('spamoor_auth');
      updateAuthUI();
      return;
    }

    /* Try to load token from sessionStorage first */
    var stored = sessionStorage.getItem('spamoor_auth');
    if (stored) {
      try {
        var data = JSON.parse(stored);
        var expiresAt = parseInt(data.expr) * 1000;
        var timeLeft = expiresAt - Date.now();

        /* If token is still valid and has more than 10 seconds left, use it */
        if (timeLeft > 10000) {
          authState.token = data.token;
          authState.user = data.user;
          authState.expiresAt = expiresAt;
          authState.isAuthenticated = true;
          authState.authDisabled = false;
          updateAuthUI();
          scheduleTokenRefresh();
          return;
        }
      } catch (e) {
        /* Invalid stored data, clear it */
        sessionStorage.removeItem('spamoor_auth');
      }
    }

    /* No valid stored token, fetch from server */
    refreshAuthToken();
  }

  // Schedule token refresh 10 seconds before expiry
  function scheduleTokenRefresh() {
    if (refreshTimer) {
      clearTimeout(refreshTimer);
      refreshTimer = null;
    }

    if (!authState.expiresAt) return;

    var refreshIn = Math.max(0, authState.expiresAt - Date.now() - 10000);
    if (refreshIn > 0) {
      refreshTimer = setTimeout(refreshAuthToken, refreshIn);
    } else {
      // Already close to expiry, refresh now
      refreshAuthToken();
    }
  }

  // Refresh auth token from server
  function refreshAuthToken() {
    fetch('/auth/token')
      .then(function(response) {
        if (response.status === 404) {
          // Auth is disabled, consider authenticated
          authState.token = null;
          authState.user = 'unauthenticated';
          authState.expiresAt = null;
          authState.isAuthenticated = true;
          authState.authDisabled = true;
          sessionStorage.removeItem('spamoor_auth');
          updateAuthUI();
          return null;
        }
        if (!response.ok) {
          throw new Error('Not authenticated');
        }
        return response.json();
      })
      .then(function(data) {
        if (data === null) return; // Auth disabled case

        if (data.token) {
          authState.token = data.token;
          authState.user = data.user;
          authState.expiresAt = parseInt(data.expr) * 1000;
          authState.isAuthenticated = true;
          authState.authDisabled = false;

          // Store in sessionStorage
          sessionStorage.setItem('spamoor_auth', JSON.stringify({
            token: data.token,
            user: data.user,
            expr: data.expr
          }));

          updateAuthUI();
          scheduleTokenRefresh();
        }
      })
      .catch(function() {
        authState.token = null;
        authState.user = null;
        authState.expiresAt = null;
        authState.isAuthenticated = false;
        authState.authDisabled = false;
        sessionStorage.removeItem('spamoor_auth');
        updateAuthUI();
      });
  }

  // Update UI based on auth state
  function updateAuthUI() {
    var loginBtn = document.getElementById('loginBtn');
    var userInfo = document.getElementById('userInfo');
    var userName = document.getElementById('userName');

    if (loginBtn && userInfo) {
      if (authState.isAuthenticated) {
        loginBtn.classList.add('d-none');
        userInfo.classList.remove('d-none');
        if (userName) {
          userName.textContent = authState.user || 'User';
        }
      } else {
        loginBtn.classList.remove('d-none');
        userInfo.classList.add('d-none');
      }
    }

    // Update page elements based on auth state
    document.querySelectorAll('[data-auth-required]').forEach(function(el) {
      if (authState.isAuthenticated) {
        el.classList.remove('d-none');
      } else {
        el.classList.add('d-none');
      }
    });

    document.querySelectorAll('[data-auth-hide]').forEach(function(el) {
      if (authState.isAuthenticated) {
        el.classList.add('d-none');
      } else {
        el.classList.remove('d-none');
      }
    });

    // Dispatch event for other scripts to listen
    window.dispatchEvent(new CustomEvent('authStateChanged', { detail: authState }));
  }

  // Get current auth token
  function getAuthToken() {
    if (!authState.isAuthenticated || !authState.token) {
      return null;
    }
    if (authState.expiresAt && Date.now() > authState.expiresAt) {
      return null;
    }
    return authState.token;
  }

  /* Check if user is authenticated */
  function isAuthenticated() {
    /* When auth is disabled, always return true */
    if (authState.authDisabled) {
      return true;
    }
    return authState.isAuthenticated && authState.token &&
           (!authState.expiresAt || Date.now() < authState.expiresAt);
  }

  // Fetch with auth token
  function authFetch(url, options) {
    options = options || {};
    options.headers = options.headers || {};

    var token = getAuthToken();
    if (token) {
      options.headers['Authorization'] = 'Bearer ' + token;
    }

    return fetch(url, options).then(function(response) {
      if (response.status === 401) {
        // Token expired or invalid, try to refresh
        refreshAuthToken();
      }
      return response;
    });
  }

  function initControls() {
    // init tooltips
    document.querySelectorAll('[data-bs-toggle="tooltip"]').forEach(initTooltip);
    cleanupTooltips();

    // init clipboard buttons
    document.querySelectorAll("[data-clipboard-text]").forEach(initCopyBtn);
    document.querySelectorAll("[data-clipboard-target]").forEach(initCopyBtn);
  }

  function initTooltip(el) {
    if($(el).data("tooltip-init"))
      return;
    //console.log("init tooltip", el);
    var idx = tooltipIdx++;
    $(el).data("tooltip-init", idx).attr("data-tooltip-idx", idx.toString());
    $(el).tooltip();
    var tooltip = bootstrap.Tooltip.getInstance(el);
    tooltipDict[idx] = {
      element: el,
      tooltip: tooltip,
    };
  }

  function cleanupTooltips() {
    Object.keys(spamoor.tooltipDict).forEach(function(idx) {
      var ref = spamoor.tooltipDict[idx];
      if(document.body.contains(ref.element)) return;
      ref.tooltip.dispose();
      delete spamoor.tooltipDict[idx];
    });
  }

  function initCopyBtn(el) {
    if($(el).data("clipboard-init"))
      return;
    $(el).data("clipboard-init", true);
    var clipboard = new ClipboardJS(el);
    clipboard.on("success", onClipboardSuccess);
    clipboard.on("error", onClipboardError);
  }

  function onClipboardSuccess(e) {
    var title = e.trigger.getAttribute("data-bs-original-title");
    var tooltip = bootstrap.Tooltip.getInstance(e.trigger);
    tooltip.setContent({ '.tooltip-inner': 'Copied!' });
    tooltip.show();
    setTimeout(function () {
      tooltip.setContent({ '.tooltip-inner': title });
    }, 1000);
  }

  function onClipboardError(e) {
    var title = e.trigger.getAttribute("data-bs-original-title");
    var tooltip = bootstrap.Tooltip.getInstance(e.trigger);
    tooltip.setContent({ '.tooltip-inner': 'Failed to Copy!' });
    tooltip.show();
    setTimeout(function () {
      tooltip.setContent({ '.tooltip-inner': title });
    }, 1000);
  }

  function updateTimers() {
    var timerEls = document.querySelectorAll("[data-timer]");
    timerEls.forEach(function(timerEl) {
      var time = timerEl.getAttribute("data-timer");
      var textEls = Array.prototype.filter.call(timerEl.querySelectorAll("*"), function(el) { return el.firstChild && el.firstChild.nodeType === 3 });
      var textEl = textEls.length ? textEls[0] : timerEl;
      
      textEl.innerText = renderRecentTime(time);
    });
  }

  function renderRecentTime(time) {
    var duration = time - Math.floor(new Date().getTime() / 1000);
    var timeStr= "";
    var absDuration = Math.abs(duration);

    if (absDuration < 1) {
      return "now";
    } else if (absDuration < 60) {
      timeStr = absDuration + " sec."
    } else if (absDuration < 60*60) {
      timeStr = (Math.floor(absDuration / 60)) + " min."
    } else if (absDuration < 24*60*60) {
      timeStr = (Math.floor(absDuration / (60 * 60))) + " hr."
    } else {
      timeStr = (Math.floor(absDuration / (60 * 60 * 24))) + " day."
    }
    if (duration < 0) {
      return timeStr + " ago";
    } else {
      return "in " + timeStr;
    }
  }

  
})()
