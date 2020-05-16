const login = new Vue({
    el: '#login',
    data: {
      // 登録内容
      mailAddress : null, 
      passwd : null      // 登録するユーザのパスワード
    },
    methods: {
      login: () => {
        const payload = {
          'mailAddress' : login.mailAddress,
          'passwd' : login.passwd
        }
        axios.post('api/login', payload)
        .then((response) => {
          alert(JSON.stringify(response));
        })
        .catch((err) => {
          alert(err)
        })
      },
      // 画面ロード時にメソッドを実行
      window:onload = function() {
      },
    }
  })