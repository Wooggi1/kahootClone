<script lang="ts">
  import svelteLogo from './assets/svelte.svg'
  import viteLogo from '/vite.svg'
  import Counter from './lib/Counter.svelte'
  import Button from './lib/Button.svelte';
  import QuizCard from './lib/QuizCard.svelte';

  let quizzes: {_id: string, name: string}[] = [];

  async function getQuizzes(){
    let response = await fetch("http://localhost:3000/api/quizzes")

    if (!response.ok){
      alert("Failed")
      return
    }

    let json = await response.json()
    quizzes = json;
  }

  let code = ""
  function connect(){
    let websocket = new WebSocket("ws://localhost:3000/ws")
    websocket.onopen = () => {
      console.log("opened connection")
      websocket.send("Hello World")
    }

    websocket.onmessage = (event) => {
      console.log(event.data)
    }
  }

  function hostQuiz(quiz) {
    let websocket = new WebSocket("ws://localhost:3000/ws")
    websocket.onopen = () => {
      console.log("opened connection")
      websocket.send(`host"${quiz.id}`)
    }

    websocket.onmessage = (event) => {
      console.log(event.data)
    }
  } 
</script>

<Button on:click={getQuizzes}>Get quizzes</Button>

<div>
  {#each quizzes as quiz }
    <QuizCard on:host={() => hostQuiz(quiz)} quiz={quiz} />
  {/each}
</div>

<input bind:value={code} class="border" type="text" placeholder="Game code" />
<Button on:click={connect}>Connect</Button>
