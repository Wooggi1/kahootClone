<script lang="ts">
  import QuizCard from "../../lib/QuizCard.svelte";
  import type { Quiz } from "../../model/quiz";

  let quizzes: Quiz[] = [];

  async function getQuizzes(): Promise<Quiz[]> {
    let response = await fetch("http://localhost:3000/api/quizzes")
    if (!response.ok) {
      alert("Failed to fetch quizzes!")
      return []
    }

    let json = await response.json()
    return json
  }

  (async function () {
    quizzes = await getQuizzes();
  })();
</script>

<div class="p-8">
  <h2 class="text-4xl font-bold">Your quizzes</h2>
  <div class="flex flex-col gap-2 mt-4">
    {#each quizzes as quiz(quiz.id)}
      <QuizCard on:host {quiz} />
    {/each}
  </div>
</div>