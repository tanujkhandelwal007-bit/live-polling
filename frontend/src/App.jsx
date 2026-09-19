import { useEffect, useState } from "react";
import "./App.css";

const API_URL = "http://localhost:8080";

function App() {
  // -------------------------
  // Authentication
  // -------------------------

  const [isRegister, setIsRegister] = useState(false);

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [token, setToken] = useState(
    localStorage.getItem("token") || ""
  );

  const [authLoading, setAuthLoading] = useState(false);

  // -------------------------
  // Poll creation
  // -------------------------

  const [question, setQuestion] = useState("");
  const [options, setOptions] = useState(["", ""]);

  // -------------------------
  // Poll display / voting
  // -------------------------

  const [poll, setPoll] = useState(null);
  const [results, setResults] = useState({});
  const [selectedOption, setSelectedOption] = useState("");

  const [loading, setLoading] = useState(false);
  const [voting, setVoting] = useState(false);
  const [closing, setClosing] = useState(false);

  const [message, setMessage] = useState("");

  // -------------------------
  // Get poll ID from URL
  // -------------------------

  const getPollIdFromURL = () => {
    const path = window.location.pathname;
    const parts = path.split("/");

    if (parts[1] === "poll" && parts[2]) {
      return parts[2];
    }

    return "";
  };

  const pollId = getPollIdFromURL();

  // -------------------------
  // Get user ID from JWT
  // -------------------------

  const getUserIdFromToken = (jwtToken) => {
    if (!jwtToken) {
      return "";
    }

    try {
      const parts = jwtToken.split(".");

      if (parts.length !== 3) {
        return "";
      }

      const payload = JSON.parse(atob(parts[1]));

      return payload.userId || "";
    } catch (error) {
      return "";
    }
  };

  const currentUserId = getUserIdFromToken(token);

  // -------------------------
  // Register / Login
  // -------------------------

  const handleAuth = async () => {
    setMessage("");

    if (!email.trim()) {
      setMessage("Please enter your email");
      return;
    }

    if (!password) {
      setMessage("Please enter your password");
      return;
    }

    if (isRegister && !name.trim()) {
      setMessage("Please enter your name");
      return;
    }

    setAuthLoading(true);

    try {
      const endpoint = isRegister
        ? "/auth/register"
        : "/auth/login";

      const body = isRegister
        ? {
          name: name.trim(),
          email: email.trim(),
          password,
        }
        : {
          email: email.trim(),
          password,
        };

      const response = await fetch(
        `${API_URL}${endpoint}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(body),
        }
      );

      const data = await response.json();

      if (!response.ok) {
        setMessage(
          data.error || "Authentication failed"
        );
        return;
      }

      if (isRegister) {
        setMessage(
          "Registration successful! Please login."
        );

        setIsRegister(false);
        setPassword("");

        return;
      }

      localStorage.setItem("token", data.token);
      setToken(data.token);

      setEmail("");
      setPassword("");

      setMessage("Login successful!");
    } catch (error) {
      setMessage("Server connection failed");
    } finally {
      setAuthLoading(false);
    }
  };

  // -------------------------
  // Logout
  // -------------------------

  const handleLogout = () => {
    localStorage.removeItem("token");
    setToken("");
    setMessage("");
  };

  // -------------------------
  // Create Poll
  // -------------------------

  const handleCreatePoll = async () => {
    setMessage("");

    if (!token) {
      setMessage("Please login first");
      return;
    }

    const cleanedQuestion = question.trim();

    const cleanedOptions = options
      .map((option) => option.trim())
      .filter((option) => option !== "");

    if (!cleanedQuestion) {
      setMessage("Please enter a poll question");
      return;
    }

    if (cleanedOptions.length < 2) {
      setMessage("Please enter at least 2 options");
      return;
    }

    setLoading(true);

    try {
      const response = await fetch(
        `${API_URL}/polls`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            question: cleanedQuestion,
            options: cleanedOptions,
          }),
        }
      );

      const data = await response.json();

      if (!response.ok) {
        if (response.status === 401) {
          localStorage.removeItem("token");
          setToken("");
          setMessage(
            "Session expired. Please login again."
          );
          return;
        }

        setMessage(
          data.error || "Failed to create poll"
        );

        return;
      }

      window.location.href =
        `/poll/${data.poll.id}`;
    } catch (error) {
      setMessage("Server connection failed");
    } finally {
      setLoading(false);
    }
  };

  // -------------------------
  // Add / Remove / Update Option
  // -------------------------

  const addOption = () => {
    setOptions([...options, ""]);
  };

  const removeOption = (index) => {
    if (options.length <= 2) {
      return;
    }

    setOptions(
      options.filter(
        (_, optionIndex) =>
          optionIndex !== index
      )
    );
  };

  const updateOption = (index, value) => {
    const updatedOptions = [...options];
    updatedOptions[index] = value;
    setOptions(updatedOptions);
  };

  // -------------------------
  // Load Poll
  // -------------------------

  useEffect(() => {
    if (!pollId) {
      return;
    }

    setLoading(true);

    fetch(`${API_URL}/polls/${pollId}`)
      .then((response) => response.json())
      .then((data) => {
        if (data.poll) {
          setPoll(data.poll);
        } else {
          setMessage("Poll not found");
        }
      })
      .catch(() => {
        setMessage("Failed to load poll");
      })
      .finally(() => {
        setLoading(false);
      });

    fetch(`${API_URL}/polls/${pollId}/results`)
      .then((response) => response.json())
      .then((data) => {
        setResults(data.results || {});
      })
      .catch(() => {
        console.log("Failed to load results");
      });

    const socket = new WebSocket(
      `ws://localhost:8080/polls/${pollId}/ws`
    );

    socket.onopen = () => {
      console.log("WebSocket connected");
    };

    socket.onmessage = (event) => {
      const update = JSON.parse(event.data);

      setResults((currentResults) => ({
        ...currentResults,
        [update.option]: update.count,
      }));
    };

    socket.onerror = () => {
      console.log("WebSocket error");
    };

    socket.onclose = () => {
      console.log("WebSocket disconnected");
    };

    return () => {
      socket.close();
    };
  }, [pollId]);

  // -------------------------
  // Vote
  // -------------------------

  const handleVote = async () => {
    if (!selectedOption) {
      setMessage("Please select an option");
      return;
    }

    if (!poll.isActive) {
      setMessage("Poll is closed");
      return;
    }

    setVoting(true);
    setMessage("");

    try {
      const response = await fetch(
        `${API_URL}/polls/${pollId}/vote`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            option: selectedOption,
          }),
        }
      );

      const data = await response.json();

      if (!response.ok) {
        setMessage(
          data.error || "Vote failed"
        );
        return;
      }

      setMessage(
        "Vote submitted successfully!"
      );

      setSelectedOption("");
    } catch (error) {
      setMessage("Server connection failed");
    } finally {
      setVoting(false);
    }
  };

  // -------------------------
  // Close Poll
  // -------------------------

  const handleClosePoll = async () => {
    if (!token) {
      setMessage("Please login first");
      return;
    }

    const confirmed = window.confirm(
      "Are you sure you want to close this poll?"
    );

    if (!confirmed) {
      return;
    }

    setClosing(true);
    setMessage("");

    try {
      const response = await fetch(
        `${API_URL}/polls/${pollId}/close`,
        {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const data = await response.json();

      if (!response.ok) {
        if (response.status === 401) {
          localStorage.removeItem("token");
          setToken("");
          setMessage(
            "Session expired. Please login again."
          );

          return;
        }

        setMessage(
          data.error || "Failed to close poll"
        );

        return;
      }

      setPoll((currentPoll) => ({
        ...currentPoll,
        isActive: false,
      }));

      setSelectedOption("");
      setMessage("Poll closed successfully!");
    } catch (error) {
      setMessage("Server connection failed");
    } finally {
      setClosing(false);
    }
  };

  // -------------------------
  // Copy Poll Link
  // -------------------------

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(
        window.location.href
      );

      setMessage("Poll link copied!");
    } catch (error) {
      setMessage("Failed to copy poll link");
    }
  };

  // -------------------------
  // Calculate total votes
  // -------------------------

  const totalVotes = poll
    ? poll.options.reduce(
      (total, option) =>
        total + Number(results[option] || 0),
      0
    )
    : 0;

  // ==================================================
  // HOME SCREEN
  // ==================================================

  if (!pollId) {
    return (
      <div className="app">
        <div className="container">
          <div className="card">

            <h1 className="title">
              Live Poll
            </h1>

            <p className="subtitle">
              Create and share polls with
              real-time results.
            </p>

            {!token ? (
              <>
                <h2 className="section-title">
                  {isRegister
                    ? "Create Account"
                    : "Welcome Back"}
                </h2>

                {isRegister && (
                  <input
                    className="input"
                    type="text"
                    placeholder="Enter your name"
                    value={name}
                    onChange={(event) =>
                      setName(event.target.value)
                    }
                  />
                )}

                <input
                  className="input"
                  type="email"
                  placeholder="Enter your email"
                  value={email}
                  onChange={(event) =>
                    setEmail(event.target.value)
                  }
                />

                <input
                  className="input"
                  type="password"
                  placeholder="Enter your password"
                  value={password}
                  onChange={(event) =>
                    setPassword(event.target.value)
                  }
                />

                <button
                  className="primary-button"
                  type="button"
                  onClick={handleAuth}
                  disabled={authLoading}
                >
                  {authLoading
                    ? "Please wait..."
                    : isRegister
                      ? "Register"
                      : "Login"}
                </button>

                <button
                  className="auth-switch"
                  type="button"
                  onClick={() => {
                    setIsRegister(!isRegister);
                    setMessage("");
                  }}
                >
                  {isRegister
                    ? "Already have an account? Login"
                    : "Don't have an account? Register"}
                </button>

                {message && (
                  <p className="message">
                    {message}
                  </p>
                )}
              </>
            ) : (
              <>
                <div className="actions">
                  <button
                    className="secondary-button"
                    type="button"
                    onClick={handleLogout}
                  >
                    Logout
                  </button>
                </div>

                <hr />

                <h2 className="section-title">
                  Create a Poll
                </h2>

                <input
                  className="input"
                  type="text"
                  placeholder="What do you want to ask?"
                  value={question}
                  onChange={(event) =>
                    setQuestion(event.target.value)
                  }
                />

                <h3>Options</h3>

                {options.map((option, index) => (
                  <div
                    className="option-row"
                    key={index}
                  >
                    <input
                      className="input"
                      type="text"
                      placeholder={`Option ${index + 1
                        }`}
                      value={option}
                      onChange={(event) =>
                        updateOption(
                          index,
                          event.target.value
                        )
                      }
                    />

                    {options.length > 2 && (
                      <button
                        className="secondary-button"
                        type="button"
                        onClick={() =>
                          removeOption(index)
                        }
                      >
                        Remove
                      </button>
                    )}
                  </div>
                ))}

                <div className="actions">
                  <button
                    className="secondary-button"
                    type="button"
                    onClick={addOption}
                  >
                    + Add Option
                  </button>
                </div>

                <br />

                <button
                  className="primary-button"
                  type="button"
                  onClick={handleCreatePoll}
                  disabled={loading}
                >
                  {loading
                    ? "Creating..."
                    : "Create Poll"}
                </button>

                {message && (
                  <p className="message">
                    {message}
                  </p>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    );
  }

  // ==================================================
  // LOADING
  // ==================================================

  if (loading || !poll) {
    return (
      <div className="app">
        <div className="container">
          <div className="card">
            <h2>
              {message || "Loading poll..."}
            </h2>
          </div>
        </div>
      </div>
    );
  }

  // -------------------------
  // Owner check
  // -------------------------

  const isOwner =
    currentUserId &&
    currentUserId === poll.createdBy;

  // ==================================================
  // POLL SCREEN
  // ==================================================

  return (
    <div className="app">
      <div className="container">
        <div className="card">

          <h1 className="title">
            Live Poll
          </h1>

          <p className="subtitle">
            Real-time voting and results
          </p>

          <div className="actions">
            <span
              className={
                poll.isActive
                  ? "status"
                  : "status closed"
              }
            >
              {poll.isActive
                ? "Live"
                : "Closed"}
            </span>
          </div>

          <h2 className="section-title">
            {poll.question}
          </h2>

          {poll.options.map((option) => (
            <label
              className="poll-option"
              key={option}
            >
              <input
                type="radio"
                name="poll-option"
                value={option}
                checked={
                  selectedOption === option
                }
                disabled={!poll.isActive}
                onChange={(event) =>
                  setSelectedOption(
                    event.target.value
                  )
                }
              />

              {option}
            </label>
          ))}

          <button
            className="primary-button"
            type="button"
            onClick={handleVote}
            disabled={
              voting || !poll.isActive
            }
          >
            {voting
              ? "Submitting..."
              : poll.isActive
                ? "Submit Vote"
                : "Voting Closed"}
          </button>

          {message && (
            <p className="message">
              {message}
            </p>
          )}

          <hr />

          <h2 className="section-title">
            Live Results
          </h2>

          <p className="subtitle">
            Total votes: {totalVotes}
          </p>

          {poll.options.map((option) => {
            const count = Number(
              results[option] || 0
            );

            const percentage =
              totalVotes === 0
                ? 0
                : Math.round(
                  (count / totalVotes) * 100
                );

            return (
              <div
                className="result-row"
                key={option}
              >
                <div className="result-label">
                  <span>{option}</span>
                  <span>
                    {count} ({percentage}%)
                  </span>
                </div>

                <div className="result-bar">
                  <div
                    className="result-fill"
                    style={{
                      width: `${percentage}%`,
                    }}
                  />
                </div>
              </div>
            );
          })}

          <hr />

          {isOwner && poll.isActive && (
            <>
              <h3>
                Poll Management
              </h3>

              <button
                className="danger-button"
                type="button"
                onClick={handleClosePoll}
                disabled={closing}
              >
                {closing
                  ? "Closing..."
                  : "Close Poll"}
              </button>

              <hr />
            </>
          )}

          <h3>
            Share this poll
          </h3>

          <button
            className="secondary-button"
            type="button"
            onClick={handleCopyLink}
          >
            Copy Poll Link
          </button>

        </div>
      </div>
    </div>
  );
}

export default App;