// 认证服务
class AuthService {
  constructor() {
    this.tokenKey = 'report_app_auth_token';
    this.userKey = 'report_app_user_info';
    // 使用相对路径，让Vite代理处理
    this.baseURL = '/api';
  }

  // 获取本地存储的token
  getToken() {
    return localStorage.getItem(this.tokenKey);
  }

  // 设置token
  setToken(token) {
    localStorage.setItem(this.tokenKey, token);
  }

  // 清除token
  clearToken() {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem(this.userKey);
  }

  // 获取用户信息
  getUser() {
    const userStr = localStorage.getItem(this.userKey);
    return userStr ? JSON.parse(userStr) : null;
  }

  // 设置用户信息
  setUser(user) {
    localStorage.setItem(this.userKey, JSON.stringify(user));
  }

  // 检查是否已登录
  isAuthenticated() {
    const token = this.getToken();
    if (!token) return false;
    
    // 检查token是否过期（简单检查）
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      const currentTime = Date.now() / 1000;
      return payload.exp > currentTime;
    } catch (error) {
      console.error('Token validation error:', error);
      return false;
    }
  }

  // 发起登录 - 支持多个提供商
  async login(provider = 'dingtalk') {
    try {
      const response = await fetch(`${this.baseURL}/auth/${provider}/login`);
      if (!response.ok) {
        throw new Error('Failed to get auth URL');
      }
      
      const data = await response.json();
      
      // 保存state和provider到localStorage
      localStorage.setItem('oauth_state', data.state);
      localStorage.setItem('oauth_provider', provider);
      
      // 跳转到授权页面
      window.location.href = data.auth_url;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  // 退出登录
  async logout() {
    try {
      const token = this.getToken();
      if (token) {
        await fetch(`${this.baseURL}/auth/logout`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });
      }
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      this.clearToken();
      // 重定向到登录页面
      window.location.href = '/login';
    }
  }

  // 处理OAuth回调（授权码）- 支持多个提供商
  async handleAuthCallback() {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');
    const state = urlParams.get('state');
    const storedState = localStorage.getItem('oauth_state');
    const provider = localStorage.getItem('oauth_provider') || 'dingtalk';
    
    if (!code) {
      throw new Error('No authorization code found');
    }
    
    if (!state || state !== storedState) {
      throw new Error('Invalid state parameter');
    }
    
    try {
      // 用授权码换取JWT token
      const response = await fetch(`${this.baseURL}/auth/${provider}/exchange`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          provider: provider,
          code: code,
          state: state
        })
      });
      
      if (!response.ok) {
        throw new Error('Failed to exchange code for token');
      }
      
      const authData = await response.json();
      
      // 保存token和用户信息
      this.setToken(authData.token);
      this.setUser(authData.user);
      
      // 清除state和provider
      localStorage.removeItem('oauth_state');
      localStorage.removeItem('oauth_provider');
      
      // 清除URL参数
      window.history.replaceState({}, document.title, window.location.pathname);
      
      return authData;
    } catch (error) {
      console.error('Auth callback error:', error);
      throw error;
    }
  }

  // 获取当前用户信息
  async getCurrentUser() {
    const token = this.getToken();
    if (!token) {
      throw new Error('No token found');
    }

    try {
      const response = await fetch(`${this.baseURL}/auth/user`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        if (response.status === 401) {
          this.clearToken();
          throw new Error('Token expired');
        }
        throw new Error('Failed to get user info');
      }

      const user = await response.json();
      this.setUser(user);
      return user;
    } catch (error) {
      console.error('Get current user error:', error);
      throw error;
    }
  }

  // 创建带认证的fetch请求
  async authenticatedFetch(url, options = {}) {
    const token = this.getToken();
    
    const authOptions = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
        ...(token ? { 'Authorization': `Bearer ${token}` } : {})
      }
    };

    try {
      const response = await fetch(url, authOptions);
      
      // 如果token过期，自动跳转到登录页面
      if (response.status === 401) {
        this.clearToken();
        window.location.href = '/login';
        throw new Error('Authentication required');
      }
      
      return response;
    } catch (error) {
      console.error('Authenticated fetch error:', error);
      throw error;
    }
  }
}

// 导出单例
export const authService = new AuthService(); 