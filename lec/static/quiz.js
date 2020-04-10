const quiz = new Vue({
    el: '#quiz',
    data: {
      ownerGroupID: 1,
    },
    methods: {
      quiz: () => {
        const payload = {'ownerGroupID' : quiz.ownerGroupID}
        axios.post('api/question/quiz', payload)
        .then((response) => {
          alert(JSON.stringify(response))
        })
        .catch((err) => {
          alert(err.response.data.error)
        })
      },
      window:onload = function() {  
        quiz.quiz();
       },
    }
  })