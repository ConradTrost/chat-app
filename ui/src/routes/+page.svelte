<script lang="ts">
  type Message = {
    clientId: string;
    message: string;
  };

  let socket: WebSocket;
  let messages: Message[] = [];
  let messageToSend = "";
  let isUsernameSet = false;
  let username = "";

  function isConnected(): boolean {
    if (socket.readyState == socket.OPEN) return true;
    else return false;
  }

  function sendChat() {
    if (isConnected()) socket.send(messageToSend);
    messageToSend = "";
  }

  function setUsername() {
    if (username.length > 3) {
      isUsernameSet = true;
      connectToSocket(username);
    } else {
      console.log("Username not long enough!");
    }
  }

  function connectToSocket(username: string) {
    username = encodeURI(username);
    socket = new WebSocket(`ws://localhost:6096/room?username=${username}`);
    socket.onopen = () => {
      console.log("Opened");
    };
    socket.onmessage = (msg: any) => {
      console.log(msg.data);
      messages = [...messages, JSON.parse(msg.data)];
    };
  }
</script>

<svelte:head>
  <title>Home</title>
  <meta name="description" content="Svelte demo app" />
</svelte:head>

{#if !isUsernameSet}
  <form on:submit={setUsername}>
    <label for="username" style="color: white">Username</label>
    <input name="username" bind:value={username} />
  </form>
{/if}
{#if isUsernameSet}
  <section>
    {#each messages as message}
      <div class={message.clientId == username ? "right" : "left"}>
        <div class="inner">
          <p>{message.message}</p>
        </div>
        {#if message.clientId !== username}
          <p class="user">{message.clientId}</p>
        {/if}
      </div>
    {/each}
  </section>

  <section>
    <form on:submit={sendChat}>
      <input bind:value={messageToSend} name="chat-input" type="text" />
      <button>Send</button>
    </form>
  </section>
{/if}

<style>
  .left,
  .right {
    display: flex;
    margin: 10px 0px;
    position: relative;
  }
  .left {
    justify-content: start;
  }
  .right {
    justify-content: end;
  }
  .left .inner,
  .right .inner {
    padding: 0 3rem;
    background-color: #fff;
    border-radius: 7%;
    width: fit-content;
    position: relative;
  }
  .right .inner {
    background-color: #9499fb;
  }
  .left .inner {
    background-color: #ececec;
  }
  .right .inner::after,
  .left .inner::after {
    content: "";
    position: absolute;
    width: 0;
    height: 0;
  }
  .right .inner::after {
    right: -8px;
    bottom: -8px;
    border-left: 20px solid transparent;
    border-right: 10px solid transparent;
    border-top: 20px solid #9499fb;
    transform: rotate(-40deg);
  }
  .left .inner::after {
    left: -8px;
    bottom: -8px;
    border-left: 10px solid transparent;
    border-right: 20px solid transparent;
    border-top: 20px solid #ececec;
    transform: rotate(40deg);
  }
  .user {
    position: absolute;
    display: block;
    left: -60px;
    max-width: 60px;
    color: white;
  }
</style>
