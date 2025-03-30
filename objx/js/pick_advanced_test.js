// 高级测试用例 - pick.js
const { pick } = require('./pick');

// 辅助函数，用于清晰显示测试结果
function printTestResult(testName, actual, expected) {
  console.log(`测试: ${testName}`);
  console.log(`期望找到${expected}个项目，实际找到: ${actual.length}`);
  
  // 简单格式化输出结果，避免太复杂的嵌套结构
  const simplifyForDisplay = (item) => {
    if (Array.isArray(item)) {
      return `Array(${item.length})`;
    }
    if (item !== null && typeof item === 'object') {
      const keys = Object.keys(item);
      if (keys.includes('name')) {
        return `{name: "${item.name}", ...}`;
      }
      if (keys.includes('id')) {
        return `{id: ${item.id}, ...}`;
      }
      if (keys.includes('title')) {
        return `{title: "${item.title}", ...}`;
      }
      return `{${keys.join(', ')}}`;
    }
    return item;
  };
  
  console.log("匹配结果:", actual.map(simplifyForDisplay));
  console.log("\n");
}

// 测试1：复杂多级嵌套结构
function testDeepNesting() {
  console.log("==== 复杂多级嵌套结构测试 ====");
  
  const data = {
    company: {
      name: "科技创新有限公司",
      founded: 2010,
      locations: ["北京", "上海", "深圳"],
      departments: [
        {
          name: "研发部",
          budget: 5000000,
          teams: [
            {
              name: "前端团队",
              members: [
                { id: 101, name: "张三", level: "高级", skills: ["JavaScript", "React", "TypeScript"] },
                { id: 102, name: "李四", level: "中级", skills: ["JavaScript", "Vue", "CSS"] }
              ],
              projects: [
                { id: "FE-001", title: "公司官网重构", priority: "高" },
                { id: "FE-002", title: "内部系统开发", priority: "中" }
              ]
            },
            {
              name: "后端团队",
              members: [
                { id: 201, name: "王五", level: "高级", skills: ["Java", "Spring", "MySQL"] },
                { id: 202, name: "赵六", level: "高级", skills: ["Go", "Docker", "Kubernetes"] }
              ],
              projects: [
                { id: "BE-001", title: "API网关", priority: "高" },
                { id: "BE-002", title: "数据处理服务", priority: "高" }
              ]
            }
          ]
        },
        {
          name: "市场部",
          budget: 3000000,
          teams: [
            {
              name: "销售团队",
              members: [
                { id: 301, name: "钱七", level: "高级", skills: ["谈判", "客户关系"] },
                { id: 302, name: "孙八", level: "初级", skills: ["市场推广", "社交媒体"] }
              ],
              projects: [
                { id: "MK-001", title: "年度营销计划", priority: "高" }
              ]
            }
          ]
        }
      ]
    }
  };
  
  // 测试 1.1: 查找所有高级开发人员
  const test1_1 = pick(data, "[level=高级]");
  printTestResult("1.1 查找所有高级开发人员", test1_1, 4);
  
  // 测试 1.2: 查找拥有JavaScript技能的成员
  const test1_2 = pick(data, "[skills*=JavaScript]");
  printTestResult("1.2 查找拥有JavaScript技能的成员", test1_2, 2);
  
  // 测试 1.3: 查找前端团队
  const test1_3 = pick(data, "[name=前端团队]");
  printTestResult("1.3 查找前端团队", test1_3, 1);
  
  // 测试 1.4: 查找所有高优先级项目
  const test1_4 = pick(data, "[priority=高]");
  printTestResult("1.4 查找所有高优先级项目", test1_4, 4);
  
  // 测试 1.5: 查找预算大于400万的部门
  const test1_5 = pick(data, "[budget=5000000]");
  printTestResult("1.5 查找预算为500万的部门", test1_5, 1);
}

// 测试2：多条件组合查询
function testMultipleConditions() {
  console.log("==== 多条件组合查询测试 ====");
  
  const data = {
    products: [
      { 
        id: "P001", 
        name: "高性能笔记本电脑", 
        category: "电子产品",
        price: 8999,
        stock: 120,
        specs: { 
          cpu: "Intel i7", 
          ram: "16GB", 
          storage: "512GB SSD",
          display: "15.6英寸 4K"
        },
        reviews: [
          { user: "user1", rating: 5, comment: "非常好用" },
          { user: "user2", rating: 4, comment: "性价比高" }
        ]
      },
      { 
        id: "P002", 
        name: "游戏笔记本", 
        category: "电子产品",
        price: 12999,
        stock: 50,
        specs: { 
          cpu: "Intel i9", 
          ram: "32GB", 
          storage: "1TB SSD",
          display: "17.3英寸 144Hz"
        },
        reviews: [
          { user: "user3", rating: 5, comment: "游戏性能极佳" },
          { user: "user4", rating: 5, comment: "散热很好" }
        ]
      },
      { 
        id: "P003", 
        name: "超薄办公本", 
        category: "电子产品",
        price: 6999,
        stock: 200,
        specs: { 
          cpu: "Intel i5", 
          ram: "8GB", 
          storage: "256GB SSD",
          display: "13.3英寸 全高清"
        },
        reviews: [
          { user: "user5", rating: 4, comment: "轻薄便携" },
          { user: "user6", rating: 3, comment: "电池续航一般" }
        ]
      },
      { 
        id: "P004", 
        name: "智能手机", 
        category: "手机",
        price: 4999,
        stock: 500,
        specs: { 
          cpu: "骁龙888", 
          ram: "12GB", 
          storage: "256GB",
          display: "6.7英寸 OLED"
        },
        reviews: [
          { user: "user7", rating: 5, comment: "相机很棒" },
          { user: "user8", rating: 4, comment: "外观设计精美" }
        ]
      }
    ],
    promotions: {
      "P001": { discount: 0.9, endDate: "2023-12-31" },
      "P002": { discount: 0.85, endDate: "2023-11-30" },
      "P004": { discount: 0.95, endDate: "2023-12-15" }
    }
  };
  
  // 测试 2.1: 查找电子产品类别且价格大于8000的产品
  const test2_1 = pick(data, "products [category=电子产品][price=8999]");
  printTestResult("2.1 查找电子产品类别且价格为8999的产品", test2_1, 1);
  
  // 测试 2.2: 查找所有带有SSD的产品
  const test2_2 = pick(data, "products [specs*=SSD]");
  printTestResult("2.2 查找所有带有SSD的产品", test2_2, 3);
  
  // 测试 2.3: 查找所有5星评价的产品
  const test2_3 = pick(data, "products reviews [rating=5]");
  printTestResult("2.3 查找所有有5星评价的评论", test2_3, 4);
  
  // 测试 2.4: 查找有促销的产品
  const productIds = Object.keys(data.promotions);
  const test2_4 = pick(data, "products").filter(product => 
    productIds.includes(product.id)
  );
  printTestResult("2.4 查找有促销的产品 (使用查询结果后处理)", test2_4, 3);
}

// 测试3：数组和对象混合嵌套
function testMixedArrayAndObjects() {
  console.log("==== 数组和对象混合嵌套测试 ====");
  
  const data = [
    {
      type: "folder",
      name: "项目文档",
      items: [
        {
          type: "document",
          name: "需求说明书.docx",
          size: 2500,
          tags: ["需求", "文档", "重要"]
        },
        {
          type: "folder",
          name: "技术规范",
          items: [
            {
              type: "document",
              name: "API设计.pdf",
              size: 3200,
              tags: ["API", "设计", "重要"]
            },
            {
              type: "document",
              name: "数据库设计.pdf",
              size: 4100,
              tags: ["数据库", "设计", "重要"]
            }
          ]
        }
      ]
    },
    {
      type: "folder",
      name: "源代码",
      items: [
        {
          type: "folder",
          name: "前端",
          items: [
            {
              type: "file",
              name: "index.js",
              size: 540,
              language: "JavaScript"
            },
            {
              type: "file",
              name: "styles.css",
              size: 320,
              language: "CSS"
            }
          ]
        },
        {
          type: "folder",
          name: "后端",
          items: [
            {
              type: "file",
              name: "server.go",
              size: 980,
              language: "Go"
            },
            {
              type: "file",
              name: "database.go",
              size: 850,
              language: "Go"
            }
          ]
        }
      ]
    }
  ];
  
  // 测试 3.1: 查找所有文档类型的文件
  const test3_1 = pick(data, "[type=document]");
  printTestResult("3.1 查找所有文档类型的文件", test3_1, 3);
  
  // 测试 3.2: 查找所有Go语言的文件
  const test3_2 = pick(data, "[language=Go]");
  printTestResult("3.2 查找所有Go语言的文件", test3_2, 2);
  
  // 测试 3.3: 查找带有"重要"标签的文档
  const test3_3 = pick(data, "[tags*=重要]");
  printTestResult("3.3 查找带有'重要'标签的文档", test3_3, 3);
  
  // 测试 3.4: 查找前端文件夹
  const test3_4 = pick(data, "items [name=前端]");
  printTestResult("3.4 查找前端文件夹", test3_4, 1);
}

// 测试4：循环引用和特殊情况
function testSpecialCases() {
  console.log("==== 循环引用和特殊情况测试 ====");
  
  // 创建带有循环引用的对象
  const cyclic = {
    name: "循环引用对象",
    value: 42,
    nestedObj: {
      name: "嵌套对象",
      tags: ["循环", "测试"]
    }
  };
  cyclic.self = cyclic; // 添加自引用
  cyclic.nestedObj.parent = cyclic; // 添加反向引用
  
  // 混合数据类型
  const mixedTypes = {
    nullValue: null,
    undefinedValue: undefined,
    numberValue: 123,
    stringValue: "test string",
    booleanValue: true,
    dateValue: new Date(),
    regexValue: /test/,
    functionValue: function() { return "test"; },
    arrayValue: [1, "2", null, undefined, { test: "nested" }],
    objectValue: { 
      a: 1, 
      b: "2",
      c: null,
      d: { nested: "value" }
    }
  };
  
  // 测试 4.1: 处理循环引用
  try {
    const test4_1 = pick(cyclic, "[name=嵌套对象]");
    printTestResult("4.1 处理循环引用", test4_1, 1);
  } catch (err) {
    console.log("测试 4.1 失败:", err.message);
  }
  
  // 测试 4.2: 混合数据类型处理
  const test4_2 = pick(mixedTypes, "[test=nested]");
  printTestResult("4.2 在混合数据类型中查找", test4_2, 1);
  
  // 非常深的嵌套
  let deepNested = { value: 1 };
  let current = deepNested;
  for (let i = 0; i < 100; i++) {
    current.next = { value: i + 2 };
    current = current.next;
  }
  
  // 测试 4.3: 非常深的嵌套
  const test4_3 = pick(deepNested, "[value=50]");
  printTestResult("4.3 在深度嵌套中查找特定值", test4_3, 1);
}

// 测试5：性能测试
function testPerformance() {
  console.log("==== 性能测试 ====");
  
  // 生成大数据集
  const generateLargeDataset = (size) => {
    const result = [];
    for (let i = 0; i < size; i++) {
      result.push({
        id: `item-${i}`,
        name: `名称 ${i}`,
        value: i % 100,
        tags: [`tag-${i % 10}`, `category-${i % 5}`],
        nested: {
          detail: `详情 ${i}`,
          code: i.toString(16),
          flags: {
            active: i % 2 === 0,
            featured: i % 5 === 0,
            special: i % 10 === 0
          }
        }
      });
    }
    return result;
  };
  
  const largeData = { items: generateLargeDataset(10000) };
  
  // 测试 5.1: 大数据集上的简单查询
  console.time("简单查询");
  const test5_1 = pick(largeData, "items [value=50]");
  console.timeEnd("简单查询");
  printTestResult("5.1 大数据集上的简单查询 (value=50)", test5_1, 1);
  
  // 测试 5.2: 大数据集上的复杂查询
  console.time("复杂查询");
  const test5_2 = pick(largeData, "items [tags*=tag-5][nested*=详情][value*=5]");
  console.timeEnd("复杂查询");
  printTestResult("5.2 大数据集上的复杂查询 (多条件)", test5_2, 1);
  
  // 测试 5.3: 特殊项查询
  console.time("特殊项查询");
  const test5_3 = pick(largeData, "items [value=0]");
  console.timeEnd("特殊项查询");
  printTestResult("5.3 查询value为0的项", test5_3, 1);
}

// 运行所有高级测试
function runAllAdvancedTests() {
  testDeepNesting();
  testMultipleConditions();
  testMixedArrayAndObjects();
  testSpecialCases();
  testPerformance();
  
  console.log("所有高级测试完成！");
}

runAllAdvancedTests(); 