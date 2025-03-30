/**
 * JavaScript version of the Pick function
 * 用于从复杂数据结构中提取特定元素的选择器工具
 */

/**
 * 转换任意值为字符串
 * @param {any} value - 要转换的值
 * @returns {string} - 转换后的字符串
 */
function toString(value) {
  if (value === null || value === undefined) {
    return "";
  }
  if (typeof value === 'string') {
    return value;
  } else if (typeof value === 'number' || typeof value === 'boolean') {
    return value.toString();
  } else if (Array.isArray(value)) {
    return value.map(v => toString(v)).join(',');
  } else if (typeof value === 'object') {
    try {
      return JSON.stringify(value);
    } catch (e) {
      return "";
    }
  }
  return String(value);
}

/**
 * 解析选择器部分，提取键名和条件
 * @param {string} selector - 选择器部分
 * @returns {object} - 包含键和条件的对象
 */
function parseSelector(selector) {
  const result = {
    key: "",
    conditions: []
  };
  
  // 提取键名，去除条件部分
  const keyPart = selector.replace(/\[.*?\]/g, '').trim();
  if (keyPart) {
    result.key = keyPart;
  }
  
  // 提取所有条件
  const conditionMatches = [...selector.matchAll(/\[(.*?)\]/g)];
  for (const match of conditionMatches) {
    const condition = match[1];
    
    // 处理不同类型的条件
    if (condition.includes('*=')) {
      const [prop, value] = condition.split('*=').map(s => s.trim());
      result.conditions.push({ 
        property: prop, 
        value: value, 
        operator: '*=' 
      });
    } else if (condition.includes('!=')) {
      const [prop, value] = condition.split('!=').map(s => s.trim());
      result.conditions.push({ 
        property: prop, 
        value: value, 
        operator: '!=' 
      });
    } else if (condition.includes('=')) {
      const [prop, value] = condition.split('=').map(s => s.trim());
      result.conditions.push({ 
        property: prop, 
        value: value, 
        operator: '==' 
      });
    }
  }
  
  return result;
}

/**
 * 检查对象是否匹配条件
 * @param {object} obj - 要检查的对象
 * @param {Array} conditions - 条件数组
 * @returns {boolean} - 是否匹配所有条件
 */
function matchesConditions(obj, conditions) {
  if (!obj || typeof obj !== 'object') {
    return false;
  }
  
  // 如果没有条件，认为匹配
  if (conditions.length === 0) {
    return true;
  }
  
  // 检查每个条件
  for (const condition of conditions) {
    const value = obj[condition.property];
    
    // 如果属性不存在
    if (value === undefined) {
      return false;
    }
    
    const strValue = toString(value);
    
    switch (condition.operator) {
      case '==':
        if (strValue !== condition.value) {
          return false;
        }
        break;
      case '*=':
        if (!strValue.includes(condition.value)) {
          return false;
        }
        break;
      case '!=':
        if (strValue === condition.value) {
          return false;
        }
        break;
    }
  }
  
  return true;
}

/**
 * 从任意对象中根据选择器规则提取元素
 * @param {any} src - 源数据对象
 * @param {string} rule - 选择器规则
 * @returns {Array} - 匹配元素数组
 */
function pick(src, rule) {
  // 处理空输入
  if (src === null || src === undefined) {
    return [];
  }
  
  if (!rule || rule.trim() === '') {
    return [src];
  }
  
  // 解析选择器为部分
  const parts = rule.split(/\s+/).filter(p => p.length > 0);
  if (parts.length === 0) {
    return [src];
  }
  
  // 解析所有选择器部分
  const parsedSelectors = parts.map(parseSelector);
  
  // 结果集
  const results = [];
  // 防止重复添加相同对象
  const seen = new Set();
  
  /**
   * 查找匹配元素的递归函数
   * @param {any} obj - 当前对象
   * @param {number} selectorIndex - 当前选择器索引
   * @param {Array} path - 当前路径
   */
  function findMatchingElements(obj, selectorIndex = 0, path = []) {
    // 基础检查
    if (!obj || selectorIndex >= parsedSelectors.length) {
      return;
    }
    
    const currentSelector = parsedSelectors[selectorIndex];
    
    // 处理对象
    if (typeof obj === 'object' && !Array.isArray(obj) && obj !== null) {
      // 检查对象是否匹配当前选择器条件
      const matchesCurrentSelector = matchesConditions(obj, currentSelector.conditions);
      
      // 如果当前对象匹配所有条件，并且是最后一个选择器，加入结果
      if (matchesCurrentSelector && selectorIndex === parsedSelectors.length - 1 && 
          (currentSelector.key === "" || path[path.length - 1] === currentSelector.key)) {
        const key = JSON.stringify(obj);
        if (!seen.has(key)) {
          seen.add(key);
          results.push(obj);
        }
      }
      
      // 遍历对象的所有属性
      for (const key in obj) {
        if (Object.prototype.hasOwnProperty.call(obj, key)) {
          const value = obj[key];
          
          // 检查键名是否匹配
          const keyMatches = currentSelector.key === "" || 
                            key.toLowerCase() === currentSelector.key.toLowerCase();
          
          // 如果键匹配且值是对象，递归检查下一个选择器
          if (keyMatches && value !== null && typeof value === 'object') {
            if (matchesCurrentSelector && selectorIndex < parsedSelectors.length - 1) {
              findMatchingElements(value, selectorIndex + 1, [...path, key]);
            }
          }
          
          // 递归检查该属性的所有子元素（重置选择器索引）
          if (value !== null && typeof value === 'object') {
            findMatchingElements(value, 0, [...path, key]);
          }
        }
      }
    } 
    // 处理数组
    else if (Array.isArray(obj)) {
      // 数组本身只能匹配空键名选择器
      const arrayMatchesSelector = currentSelector.key === "" && 
                                  matchesConditions(obj, currentSelector.conditions);
      
      if (arrayMatchesSelector && selectorIndex === parsedSelectors.length - 1) {
        const key = JSON.stringify(obj);
        if (!seen.has(key)) {
          seen.add(key);
          results.push(obj);
        }
      }
      
      // 遍历数组的所有元素
      for (let i = 0; i < obj.length; i++) {
        const item = obj[i];
        
        // 检查数组元素是否匹配当前选择器条件
        if (item !== null && typeof item === 'object') {
          const itemMatchesConditions = matchesConditions(item, currentSelector.conditions);
          
          // 如果元素匹配条件且是最后一个选择器，加入结果
          if (itemMatchesConditions && selectorIndex === parsedSelectors.length - 1) {
            const key = JSON.stringify(item);
            if (!seen.has(key)) {
              seen.add(key);
              results.push(item);
            }
          }
          
          // 如果元素匹配条件且不是最后一个选择器，继续匹配下一个选择器
          if (itemMatchesConditions && selectorIndex < parsedSelectors.length - 1) {
            findMatchingElements(item, selectorIndex + 1, [...path, i]);
          }
        }
        
        // 递归检查该元素（重置选择器索引）
        if (item !== null && typeof item === 'object') {
          findMatchingElements(item, 0, [...path, i]);
        }
      }
    }
  }
  
  // 开始查找
  findMatchingElements(src);
  
  return results;
}

module.exports = {
  pick
}; 