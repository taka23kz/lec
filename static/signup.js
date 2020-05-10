const signup = new Vue({
    el: '#signup',
    data: {
      // 登録内容
      userName : null,      // 登録するユーザ名
      password : null,      // 登録するユーザのパスワード
      password2 : null,      // 登録するユーザのパスワード(確認用)
    },
    methods: {
      signup: () => {
        const payload = {
          'userName' : signup.userName,
          'password' : signup.password,
          'password2' : signup.password2
        }
        axios.post('api/signup/signup', payload)
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