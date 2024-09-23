<script lang="ts">
  import Button from './lib/Button.svelte';
  import QuizCard from './lib/QuizCard.svelte';
  import { NetService } from './service/net';
  import type { Quiz } from './model/quiz';

  let quizzes: {_id: string, name: string}[] = [];

  let netService = new NetService();
  netService.connect

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
  let msg = ""

  function connect(){
    netService.sendPacket({
      id: 0,
      code: "1234",
      name: "coolname1234"
    })
  }

  function hostQuiz(quiz: Quiz) {
    netService.sendPacket({
      id: 1,
      quizId: quiz.id
    })
  } 
</script>

<Button on:click={getQuizzes}>Get quizzes</Button>
Message: {msg}

<div>
  {#each quizzes as quiz }
    <QuizCard on:host={() => hostQuiz(quiz)} quiz={quiz} />
  {/each}
</div>

<input bind:value={code} class="border" type="text" placeholder="Game code" />
<Button on:click={connect}>Connect</Button>
